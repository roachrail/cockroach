- Feature Name: Tenant capabilities
- Status: completed
- Start Date: 2022-08-10
- Authors: knz
- RFC PR: [#85954](https://github.com/cockroachdb/cockroach/pull/85954)
- Cockroach Issue: TBD

# Summary

This RFC proposes to extend tenant records with **capabilities**,
which are *run-time configurable* attributes of a tenant that determine
whether that tenant can perform certain operations outside of its
tenant keyspace.

An example capability would be the ability of a tenant to manually
relocate ranges to particular nodes. (Currently secondary tenants do
not have this capability.) Capabilities are granted by the
system/admin tenant and cannot be controlled by secondary tenants
themselves.

In essence, the introduction of capabilities decouples *resource usage
and measurement isolation*, which we will continue to support via
tenant boundaries, and *authorization for KV operations*, which we
will now start to implement via grantable capabilities.

Capabilities will help the following use cases:

- they will power a more flexible, less complex roadmap towards a
  multi-tenant DB Console.
- they will reduce the need to grant all SREs access to the SQL `admin`
  role on the system/admin tenant, i.e. enable more fine-grained, less
  vulnerable authorization rules for CockroachCloud SREs.
- they will reduce potential disruption to SRE/DBA monitoring tools
  caused by the introduction of multi-tenancy in new deployments, by
  making the "application" tenant look & feel more like a system/admin
  tenant in most regards.
- it will enable a key mechanism needed to migrate previous
  single-tenant deployment to multi-tenancy.

We will introduce capabilities using a new field in the `info` column
of `system.tenants`. Capabilities will be subsequently checked by the
KV side of the tenant-KV RPCs upon each request (i.e. we will not
trust tenant SQL servers to enforce their own capabilities, unlike SQL
privileges).

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Summary](#summary)
- [Motivation](#motivation)
    - [Fine-grained authorization for cluster operations](#fine-grained-authorization-for-cluster-operations)
    - [Separating resource isolation and control isolation](#separating-resource-isolation-and-control-isolation)
    - [Enabling our multi-tenancy roadmap](#enabling-our-multi-tenancy-roadmap)
- [Technical design](#technical-design)
    - [List of capabilities](#list-of-capabilities)
    - [Mechanisms](#mechanisms)
        - [Capability checks for features already using networked KV APIs](#capability-checks-for-features-already-using-networked-kv-apis)
        - [Capability checks for features not already using networked KV APIs](#capability-checks-for-features-not-already-using-networked-kv-apis)
        - [Administrating capabilities](#administrating-capabilities)
        - [Back-end architecture](#back-end-architecture)
    - [Usage in CockroachCloud](#usage-in-cockroachcloud)
        - [Advanced security: smaller-scope SRE access](#advanced-security-smaller-scope-sre-access)
    - [Usage when setting up new self-hosted clusters](#usage-when-setting-up-new-self-hosted-clusters)
    - [Usage when migrating single-tenant clusters into multi-tenancy](#usage-when-migrating-single-tenant-clusters-into-multi-tenancy)
    - [Drawbacks](#drawbacks)
    - [Rationale and Alternatives](#rationale-and-alternatives)
- [Explain it to folk outside of your team](#explain-it-to-folk-outside-of-your-team)
- [Unresolved questions](#unresolved-questions)
- [Footnotes](#footnotes)

<!-- markdown-toc end -->


# Motivation

## Fine-grained authorization for cluster operations

Today, the capabilities needed to administrate a storage cluster are
all granted to the "system tenant"—which is fine—, but also *only*
granted to the system tenant. We find this restriction lacking in
multiple areas.

- we plan to introduce multi-tenancy to all deployments, and
  segregating application workload and observability scope to a new
  secondary tenant, hereafter called "application tenant". In this
  world, we would like to offer self-hosted users, who are still
  choosing self-hosted for the privilege to administrate their cluster
  how they see fit, with the option to access certain "administrative
  operations" from their application tenant. Here are some example
  such operations:

  - the ability to query cluster-wide timeseries, such as node-node
    latencies.
  - configuring certain cluster setting that have cluster-wide impact,
    such as the KV rebalance rate.
  - manually placing replicas/leaseholders on specific nodes via the
    ALTER RANGE statement.
  - the ability to pre-split ranges in SQL indexes.
  - the ability to set up zone configs for meta ranges.

- we have also, through the past few years, learned from experience
  that our users are dissatisfied with coarse-grained authz rules:
  when we build two different mechanisms, our users reliably tell us
  they want *separate privileges* to access them, so that they can
  grant one but not the other.  Likewise, in the context of tenant
  capabilities, we expect that we will not be served well by a single
  binary capability "has system tenant capabilities" / "does not have
  system tenant capabilities". Instead, we want *discrete* capabilites
  that grant access/control to separate features, and can be
  granted to / revoked from tenants separately.

## Separating resource isolation and control isolation

Another angle to motivate capabilities is our growing collective
realization that we may have overly conflated the matters of *resource
and observability isolation*, which help us attribute data /
processing usage to specific entities (essentially, an *accounting*
matter); with the matters of *authorization*, which determine who
has access to what.

Until now, the two were conflated as follows: system tenant did get
permissive authorization but not accounting, and secondary tenants got
accounting and restrictive authorization.

But this is not the only way to do things. Here it's interesting to
observe a parallel to *processes* in operating systems: processes
provide resource isolation, but can be granted *permission* from the
outside to operate on system-wide resources. The same running process
can also gain and lose these privileges through its execution.

We would like to evolve CockroachDB in this direction.

## Enabling our multi-tenancy roadmap

Concurrently with the above, we have observed mounting pressure for
_some_ solution for the following two problems:

- our observability roadmap for DB Console would like to create a shared
  engineering plan for two different UX stories:

  - in Cockroach Cloud, the Console components should be restricted
    to only observe the resources specific to the tenant they are pointed
    to.
  - in Self-hosted CockroachDB, we want to offer the *option* for
    backward-compatibility for a DB Console pointed to an application
    (secondary) tenant to also observe certain items outside of that
    tenant's boundary. For example, we want the option to make the hot
    range heatmap include meta ranges.

  Without a tenant capability model, the obs inf team would need to
  create new ad-hoc mechanisms to include certain data in one
  deployment, and other data in the other. It would also be hard(er)
  to test both, as the engineers would need to change deployment
  styles to trigger different behaviors.

  With a capability model, that determination is shared and can be
  changed at run-time for testing purposes.

- when we consider upgrading single-tenant clusters to multi-tenancy,
  we want to implement a migration that logically moves the
  previously-created application data outside of the system/admin
  tenant to become a regular secondary tenant.  We also would like
  this migration to happen online, and avoid downtime.

  The plan we're considering for this is a *pivot migration*: the
  workload is running on tenant ID 1, with the capabilities otherwise
  available to system/admin tenant. Then we create another secondary
  tenant.  Then we grant the system/admin tenant capabilities to that
  new secondary tenant. Then we revoke the system/admin tenant
  capabilities from the tenant with ID 1.

  The result is that we still have a system/admin tenant, but it's not
  any more the tenant with ID 1.

# Technical design

We will introduce tenant capabilities via a new field in the `info`
column of `system.tenants`. For this we will use richer data
structures than simple booleans.

When performing an operation in SQL that is subject to tenant
capabilities, we will then change the condition from `if
currentTenantID = SystemTenantID`, to a different mechanism that
inspects the current tenant capabilities. This is described below.

## List of capabilities

We envision the list of capabilities to include (but not be limited to):

| Capability                             | Description                                                                           |
|----------------------------------------|---------------------------------------------------------------------------------------|
| `can_view_tsdb_metrics`                | Read timeseries for the KV-level metrics. [^ts]                                       |
| `can_admin_relocate_range`             | Move replicas or leasholders of the tenant's own ranges to specific nodes.            |
| `can_admin_scatter`                    | Scatter replicas or leaseholders for the specified tenant's ranges to specific nodes. |
| `can_admin_split`, `can_admin_unsplit` | Split/unsplit the tenant's own ranges manually.                                       |
| `can_check_consistency`                | Trigger a replica consistency check manually on the tenant's own ranges.              |
| `can_use_nodelocal_storage`            | Use the `nodelocal://` storage in bulk I/O statements.                                |
| `can_view_node_info`                   | Access the KV node details.                                                           |
| `can_view_storage_metrics`             | Access the time series of the storage cluster.                                        |
| `exempt_from_rate_limiting`            | Perform KV operations without rate and burst limiting.                                |
| `span_config_bounds`                   | Set zone configs with parameters within specified bounds.                             |

[^ts]: See the separate proposal on [Timeseries access for tenants](https://github.com/cockroachdb/cockroach/pull/86524).

## Mechanisms

### Capability checks for features already using networked KV APIs

Generally, the rule will be that a capability check must be performed *at least*
on the KV side of the SQL/KV boundary. Our security model will be that tenants
should not be the sole entity responsible for enforcing tenant capabilities
(we do not want a rogue tenant process to elevate their own capabilities.)

However, the check can *also* be performed in the SQL layer, to improve UX.

Let's take the example of ALTER RANGE RELOCATE.

Currently this code does this:

1. check if tenant ID == 1, otherwise report an error to the user
   preemptively.
2. send an AdminRelocate request to the KV Batch interface.
3. on the KV side of the tenant connector, AdminRelocate is refused if
   the origin tenant ID is not 1.

Notice how step 1 is optional: if it wasn't there, we would still get
an error at step 3 if the tenant was insufficiently privileged—albeit
perhaps a less informative error.

We will continue to use this overall structure, and evolve it as follows:

1. check if the current tenant currently has the `can_admin_relocate_range` capability
   for the target range(s). If not, report an error to the user.
2. send an AdminRelocate request to the KV Batch interface, as before.
3. on the KV side, AdminRelocate is refused if the origin tenant
   does not have the `can_admin_relocate_range` capability for the target range.

### Capability checks for features not already using networked KV APIs

We also have features which operate as follows:

1. check if tenant ID == 1, if not report an error to the user.
2. perform a privileged operation.

For example, accessing the KV node details in DB console.

With the capability model, we can't just keep this structure
and change step 1 to just check the capabilities. **We also need
to redirect the operation to the appropriate place in the system**.

If, say, the current system/admin tenant has ID 1,
and we grant `can_view_node_info` to tenant with ID 2,
and an operator tries to access the API endpoint to retrieve node details,
we can't keep the current code that just retrieves the data from gossip. There may
be no gossip instance locally to work with.

Instead, we must **create a new KV/SQL networked API** which channels
the operation from tenant 2's server process, through the KV layer, to
become an operation in tenant 1's server.

### Administrating capabilities

We would introduce new SQL syntax: `ALTER TENANT ... GRANT/REVOKE CAPABILITY`.

Note that the permission to change capabilities for tenants will *not*
be itself a capability. It will be solely linked to the current SQL
user's ability to write to the *current tenant's* `system.tenants`
table.

### Back-end architecture

There will be a cache on every KV node of current capabilities for
tenants that are performing KV requests. The capabilities will be
refreshed using a rangefeed.

We do not provide a synchronisation mechanism at this point, meaning
that the `ALTER TENANT GRANT/REVOKE CAPABILITY` statement has only
eventual consistency.

## Usage in CockroachCloud

When setting up new CC Dedicated clusters in a multi-tenant configuration,
we want to segregate the user's activity in such a way that they cannot
disrupt the deployment and break Cockroach Labs' own SLOs.

However, since the customers are paying for the privilege to use "all" the
hardware resources, we should not limit actions that merely incur
CPU/disk/memory usage that we don't know how to measure at a fine
grain, like we need to do in CC Serverless.

So we will set up new clusters by granting the following capabilities
to the application tenant:

- **with** the `can_view_storage_metrics` capability.
- **with** the `can_admin_relocate_range` capability.
- **with** the `can_admin_scatter`, `can_admin_split`, `can_admin_unsplit` capability.
- **with** the `can_view_node_info` capability on all the system's ranges.
- **without** capabilities on meta/liveness/etc (non-tenant ranges).
- **with** the `exempt_from_rate_limiting` capability.
- **with** suitable `span_config_bounds` that let the users apply
  custom zone configs to their own ranges.

### Advanced security: smaller-scope SRE access

For many maintainance tasks, a SRE does not need to have complete
control over the cluster. We could seek to set up an environment where
each task is carried out according to the [principle of least
privilege](https://en.wikipedia.org/wiki/Principle_of_least_privilege).

With tenant capabilities, we can achieve this as follows.

When a SRE needs to carry out a task A (e.g. configure a specific
cluster setting, extract some data, set up a backup schedule, inspect
jobs), which in turn needs capability C, we could do as follows:

- SRE requests from orchestration access to a secondary tenant with
  capability C, within a specific time frame.
- orchestration sets up a tenant with just that capability, and
  performs audit logging about this fact.
- SRE gets access to just that secondary tenant.
- SRE performs task A.
- orchestration enforces the lifetime of the access by force-deleting
  the secondary tenant when the access expires.

With this model, at no point does the SRE need access to the
system/admin tenant with all capabilities enabled.

## Usage when setting up new self-hosted clusters

We will use capabilities to decouple the introduction of multi-tenancy
and that of capability separation for self-hosted customers.

- In a first phase, we will implement the creation of an application
  tenant with *all* the capabilities otherwise reserved to the system
  tenant. In that phase, the app workloads start running in a
  secondary tenant but at this point there's no difference in
  capability with the system tenant.

  This phase makes it possible to introduce multi-tenancy for resource
  isolation without requiring a change in docs and runbooks.

- In a second phase, for new clusters we will also start reducing
  the capability set auto-granted to the application tenant.

  This second phase will be accompanied by an incremental docs project
  to cover which product areas are evolving for each of the
  capabilities removed from the app tenant.


Separately, we will also teach the users about the capability model,
and how to revoke capabilities that were auto-granted to the application
tenant, in case they are interested in reducing the control scope of their
app tenant (this will likely be desirable for Enterprise deployments).

## Usage when migrating single-tenant clusters into multi-tenancy

Overall, the migration process for bringing a single-tenant cluster
into multi-tenancy will look like this:

1. create a new secondary tenant, with all the capabilities
   otherwise reserved to the system/admin tenant.
2. remove the advanced capabilities from the tenant with ID 1
   (previous system tenant), where all the application data is stored.

The particulars of this migration will be covered by a separate RFC.

## Drawbacks

This introduces new complexity in our multi-tenant model.

## Rationale and Alternatives

The main alternative is "do nothing at all".

We are dissatisfied with this alternative because it forces us to
complete our entire roadmap for secondary tenant features before we
can even start offering multi-tenancy to existing Self-hosted or CC
Dedicated users. Tenant capabilities allow us to decouple the two
roadmaps.

# Explain it to folk outside of your team

In multi-tenant CockroachDB, we are changing how it is decided whether
a tenant can observe or administrate certain parts of the entire
(cross-tenant) storage cluster.

Previously, whole-cluster operation was limited to a special tenant
called "system tenant". An operator could only perform these operations
if it was connected to that tenant's SQL interface.

Instead, we are now going to decide whether these operations
are possible based on a new **capability model**, akin to privileges.
Capabilities can be granted to or revoked from tenants through
their lifecycle, not only when tenants are created.

The capabilities constrain which features are available within the
tenant SQL interface, even to SQL users with the `admin` role.

Some examples:

| Feature                                    | Condition for use, prior         | Condition for use, new                                                    |
|--------------------------------------------|----------------------------------|---------------------------------------------------------------------------|
| `ALTER RANGE RELOCATE`                     | `admin` SQL role + system tenant | `admin` SQL role + tenant capability `can_admin_relocate_range`           |
| View cluster-wide timeseries in DB Console | Use DB Console via system tenant | Use DB Console with any tenant with capability `can_view_storage_metrics` |

Capabilities are also granted from "outside" of secondary tenants, and
are not available by default (to any SQL user, even `admin`) unless
granted.

# Unresolved questions

N/A

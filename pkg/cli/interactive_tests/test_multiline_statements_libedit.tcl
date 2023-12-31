#! /usr/bin/env expect -f

source [file join [file dirname $argv0] common.tcl]

start_server $argv

# force use of libedit.
# NB: these tests will disappear when libedit goes away.
set env(COCKROACH_SQL_FORCE_LIBEDIT) "true"

spawn $argv sql
eexpect "defaultdb>"

start_test "Test that a multi-line entry can be recalled escaped."
send "select 'foo\r"
eexpect " ->"
send "bar';\r"
eexpect "1 row"
eexpect "defaultdb>"

# Send up-arrow.
send "\033\[A"
eexpect "select 'foo"
eexpect "bar';"
send "\r"
eexpect foo
eexpect bar
eexpect "defaultdb>"

send "select 1,\r"
eexpect " ->"
send "2, 3\r"
eexpect " ->"
end_test

start_test "Test that \p does what it says."
send "\\p\r"
eexpect "select 1,"
eexpect "2, 3"
eexpect " ->"
end_test

start_test "Test finishing the multi-line statement."
send ";\r"
eexpect "1 row"
eexpect "defaultdb>"

# Send up-arrow.
send "\033\[A"
eexpect "select 1,"
eexpect "2, 3"
eexpect ";"
end_test

start_test "Test that \r does what it says."
# backspace to erase the semicolon
send "\010"
# newline to get a prompt
send "\r"
eexpect " ->"
# Now send \r followed by a carriage return.
send "\\r\r"
eexpect "defaultdb>"
end_test

start_test "Test that Ctrl+C after the first line merely cancels the statement and presents the prompt."
send "\r"
eexpect "defaultdb>"
send "select\r"
eexpect " ->"
interrupt
eexpect "defaultdb>"
end_test

start_test "Test that \p does what it says."
send "select\r"
eexpect " ->"
send "\\p\r"
eexpect "select\r\n*->"
interrupt
eexpect "defaultdb>"
end_test

start_test "Test that a dangling table creation can be committed, and that other non-DDL, non-DML statements can be issued in the same txn. (#15283)"
send "create database if not exists t;"
send "drop table if exists t.blih;"
send "create table if not exists t.kv(k int primary key, v int);\r"
eexpect "CREATE TABLE"
eexpect "defaultdb>"
send "begin; create table t.blih(x INT REFERENCES t.kv(k));\r\r"
eexpect "CREATE TABLE"
eexpect "defaultdb"
eexpect OPEN

send "show all cluster settings;\r"
eexpect "rows"
eexpect "defaultdb"
eexpect OPEN

send "commit;\r"
eexpect COMMIT
eexpect "defaultdb>"
end_test

send "quit\r"
eexpect eof

# we force TERM to xterm, otherwise we can't
# test bracketed paste below.
set env(TERM) xterm

spawn $argv sql
eexpect "defaultdb>"

start_test "Test that a multi-line bracketed paste is handled properly."
send "\033\[200~"
send "\\set display_format csv\r\n"
send "values (1,'a'), (2,'b'), (3,'c');\r\n"
send "\033\[201~\r\n"
eexpect "1,a"
eexpect "2,b"
eexpect "3,c"
eexpect "defaultdb>"
end_test



send_eof
eexpect eof

stop_server $argv

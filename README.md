# mycnf

Read a .my.cnf file and extract credentials and other details to connect to a MySQL database.

`func ReadMyCnf` reads a .my.cnf (*~/.my.cnf* by default), searches for the named "profile" in that file and returns a DSN (string) suitable for use in `sql.Open()`.

This is based on a gist from https://gist.github.com/nickcarenza/d847ec24455e70a8609b6602ed528133 and modified to accept default values passed in a structure from the calling function.

# This is a comment at the beginning of the file

9,15H
$g
9,15d

4P
5p
s/foo/bar/g
s/t//g
s/more/MORE/1
s/[oO]/0/2
#5,s/[iI]/1/2
# n
3a Append text\
with another line\
and one more.\

s/[Ll]ine/form/g
s/s/5/g
#$/q/-1
/O[RF]/s/O/0h/g

6iInsert text\
with another line\
and one more.\

1,7i: 
1,7=

#6D

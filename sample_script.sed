4P
5p
s/foo/bar/g
s/t//g
s/more/MORE/1
s/[oO]/0/2
5,s/[iI]/1/2
# n
# d
3a/ Append text\
with another line\
and one more.
s/[Ll]ine/form/g
s/s/5/g
#$/q/-1

6i/Insert text\
with another line\
and one more. 

i/: 
=

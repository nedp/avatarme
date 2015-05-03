# avatarme
Simple avatar-making web server, for Go coding practice.

Currently (2015-05-03) hosted at:
http://avatarme.herokuapp.com/username_goes_here.png

Supports query string parameters:

 Param | Default | Description
------:| -------:|:-----------
  s    |     130 | Desired image size (width and height) in pixels. Only approximate, as a simplification to eliminate rounding errors.
  n    |      11 | Number of 'blocks' (width and height). Blocks are squares which are either white or coloured.
  b    |       1 | Number of white blocks to use as a border. Border blocks don't count towards the number specified by n.

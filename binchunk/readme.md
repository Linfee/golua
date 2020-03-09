### lua执行原理

lua脚本并不是有lua解释器解释执行，lua内置了虚拟机，会先将lua脚本编译为自己吗，然后再有虚拟机去执行

lua字节码需要一个载体，这个载体就是二进制chunk，类似于java的class

### chunk


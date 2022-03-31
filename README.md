# sstable

## [Data structure desc](<http://mezhov.com/2013/09/sstable-lsm-tree.html>)

## Implementation

### Table format
```
[header]
[block][block][block]
[index][index][index]
```

### Table spec
```
header - 8bit structure

[:4] - file length // 4bit uint32 value
[4:8] - first index block offset // 4bit uint32 value
```
```
block - row which stores key:value
each block contains following fields

[:4] - key length // 4bit uint32 value
[4:8] - value length // 4bit uint32 value
[8:kl+8] - stored key
[kl+8] - stored value

so it looks like this im memory representation
for "1":"22" value
[<0,0,0,1>,<0,0,0,2>,<1>,<2,2>]
```
```
index - each row contains block offset in file,
key length and particular key

[:4] - block offset // 4bit uint32 value
[4:8] - key length // 4bit uint32 value
[8:kl+8] - key
```

## How it works
### Write
SSTable receives a map[string]string, sorts it,
writes to file with calculated index for each block(key:val row)

### Read
When .sst loads to SSTable struct, it checks for header to get info about where to find index block,
then it loads index to memory and stores it in struct.
After this, Search method will check this index for key
and if it founds it, table goes to particular file offset, reads block and returns its value.
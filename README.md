# extract-vmlinux-v2
golang implementation of [extract-vmlinux](https://github.com/torvalds/linux/blob/master/scripts/extract-vmlinux) script

## Installing

```
go install github.com/Caesurus/extract-vmlinux-v2@latest
```

TODO: Alternative methods to come...

---
## Compressions

### Supported
At the moment I have only found good native support for `gzip`/`bzip`/`lzma`/`lz4`/`xz`.

### Not Supported
I haven't found good support for `zstd`/`lzop`.

---
## Why ???
The extract-vmlinux script that is provided in the linux repo only really caters to x86. It checks whether the extracted output is an ELF file, but on may embedded systems, the extracted vmlinux file is not an elf file recognized by `readelf`.
```bash
check_vmlinux()
{
	# Use readelf to check if it's a valid ELF
	# TODO: find a better to way to check that it's really vmlinux
	#       and not just an elf
	readelf -h $1 > /dev/null 2>&1 || return 1

	cat $1
	exit 0
}
```
Whats the alternative? The general consensus is to just run [binwalk](https://www.refirmlabs.com/binwalk/) on the `zImage`/`vmlinuz`/`bzImage` etc... and that definitely works. However it extracts anything it thinks is a file, which is by design, but this can produce extra cruft, and so it's left to the user to sift through the resulting files to find which one contains the vmlinux image you actually want.


This same thing could be achived by modifying the orginal bash script, and mailing list evidence shows plenty have done this and submitted patches that have promptly been ignored. I thought it would be fun to implement this in `go` with the thought that if it was all native `go` code, the user could use this on any OS/Architecture that `go` supports by just doing a `go get`. So now you should be able to easily extract `vmlinux` files on Windows (why you'd want to do so is beyond me, but it should work)


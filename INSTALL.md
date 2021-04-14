
Installation
============

1. Download a release from `https://github.com/runeimp/chronus/releases`
2. For archives or packages:
	* For Zip archives (`.zip`) I recommend:
		1. `unzip chronus_1.0.0_windows_x86_64.zip`
		2. `cd chronus_1.0.0_windows_x86_64`
		3. `copy chronus.exe YOUR_PREFERED_PATH`
	* For Tar-Gzip archives (`.tar.gz`) I recommend:
		1. `tar xfz chronus_1.0.0_freebsd_x86_64.tar.gz`
		2. `cd chronus_1.0.0_freebsd_x86_64`
		3. `cp chronus YOUR_PREFERED_PATH` or `ln -s chronus YOUR_PREFERED_PATH/chronus` if your source folder is not going to get deleted later.
	* For the Debian packages (`.deb`) I recommend:
		* `sudo dpkg -i /path/to/chronus_1.0.0_linux_x86_64.deb` to install and upgrade. There are no dependencies so dpkg is fine for this task.
	* For the RedHat Package Manager (`.rpm`) I recommend:
		* `rpm -i /path/to/chronus_1.0.0_linux_x86_64.rpm` to install
		* `rpm -U /path/to/chronus_1.0.0_linux_x86_64.rpm` to upgrade
3. Copy the binary to a directory in your PATH and make sure it is executable for your system. It _should_ be ready to go by default.


Finding YOUR_PREFERED_PATH (Windows)
------------------------------------

On Windows systems you can review your path list by:

1. Open a Command Prompt
2. Type `ECHO "%PATH:;="&ECHO "%"`


Finding YOUR_PREFERED_PATH (POSIX)
----------------------------------

On POSIX systems (\*BSD, Linux, macOS, Solaris, UNIX, etc.) you can review your path list with

```text
$ echo $PATH | tr : "\n"
/Users/runeimp/.local/bin
/Users/runeimp/bin
/Users/runeimp/dev/lang/go/bin
/opt/local/bin
/usr/local/bin
/usr/local/sbin
/usr/bin
/usr/sbin
/bin
/sbin
```

If you use the command `echo $PATH | tr : "\n"` you should see a list similar to the example with `runeimp` replaced with whatever your username is on the system. Directory search priority is from top to bottom when you type in a command.


### Which Paths to Avoid

On any system there are certain paths you should not install apps into. On POSIX systems any path that ends with `sbin` means system binary. **NEVER** manually install apps into an `sbin` path. Following is a general order of preference when installing apps on a system for all users:

* `/opt/local/bin` - This is generally best as `opt` paths are intended for _option software_ when it is present.
	* This is where most package managers will install binaries. If it doesn't exist you should probably install something wonderful like [`just`][] or [`kitty`][] with your package manager to kickstart the path.
* `/usr/local/opt` - This is a good path to use if it exists in your list and `/opt/local/bin` is not available.
	* This is where many package managers will install binaries. If it doesn't exist you should probably install something wonderful like [`git`][] or [`wget`][] with your package manager to kickstart the path.
* `/usr/local/bin` - This is a good path to use if neither `opt` path is available. Almost always available on all POSIX systems.
	* This is where most package managers will install binaries. If it doesn't exist you should probably install something wonderful like [`jq`][] or [`htop`][] with your package manager to kickstart the path.
* `/usr/bin` - This is acceptable if none of the above paths are available. Though it should be avoided if possible.
	* If this seems to be your only option **definitely** try to install something with your package manager if you haven't already. Doing so should create one of the above paths which are _much_ preferred over `/usr/bin`.
* `/bin` - While you could _technically_ install binaries to this path it's a very bad idea. System updates expect full control over this path and may remove or overwrite binaries you install there without warning before hand or notice after the fact.





[`git`]: https://git-scm.com/
[`htop`]: https://hisham.hm/htop/
[`jq`]: https://stedolan.github.io/jq/
[`just`]: https://github.com/casey/just
[`kitty`]: https://sw.kovidgoyal.net/kitty/
[`wget`]: https://www.gnu.org/software/wget/


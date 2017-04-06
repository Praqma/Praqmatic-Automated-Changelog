# Installing PAC on Windows

1. Download and install [Ruby](https://dl.bintray.com/oneclick/rubyinstaller/rubyinstaller-2.3.3-x64.exe), let the installer add Ruby to PATH when given the option.
2. Download [DevKit](https://dl.bintray.com/oneclick/rubyinstaller/DevKit-mingw64-64-4.7.2-20130224-1432-sfx.exe) and install it following [the DevKit quick-start](https://github.com/oneclick/rubyinstaller/wiki/Development-Kit#quick-start).
3. Download and install [cmake](http://www.cmake.org/files/v3.2/cmake-3.2.3-win32-x86.zip) and [pkg-config](http://iweb.dl.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win64/Personal%20Builds/ray_linn/64bit-libraries/pkg-config/pkg-config-0.26.7z) following instructions detailed in [this SO answer](http://stackoverflow.com/a/31254515).
4. Download [wkhtmltopdf](https://downloads.wkhtmltopdf.org/0.12/0.12.4/wkhtmltox-0.12.4_msvc2015-win64.exe) and install it to a path that does not contain any spaces, for example `C:\tools\wkhtmltopdf`.
5. Add an [environment variable](https://kb.wisc.edu/cae/page.php?id=24500) with the name `wkhtmltopdf` and the full path to the executable as value, for example `C:\tools\wkhtmltopdf\bin\wkhtmltopdf.exe`.
6. Install `bundler` by running `gem install bundler`
7. Clone the pac repository to your local machine: `git clone https://github.com/Praqma/Praqmatic-Automated-Changelog.git`.
8. Optionally check-out the latest tag or a specific release tag if you don't want bleeding edge.
9. Change directory to Praqmatic-Automated-Changelog (the git clone) and run the command `bundle install` to install all the used Ruby Gems.

## Testing the installation

Open `default_settings.yml` and replace the `:task_systems:` block with the following:

```
:task_systems:
  - :name: none
    :regex:
      - { pattern: '/(CS-\d+)/i', label: none }
    :delimiter: '/,|\s/'
```

Set the repository location to the `misc` folder in the PAC project. 

`:repo_location: '..\Praqmatic-Automated-Changelog\misc'`

In the above mentioned folder, create a sample repository containing a commit with the message `CS-1`, and a tag with the name `v0.1`.
You can execute the following code snippet in the `misc` folder to quickly set up a small test repository:
```
git init
git commit --allow-empty -m "Initial commit"
git tag -a v0.1 -m "version 0.1"
git commit --allow-empty -m "CS-1"
```

Run `ruby pac.rb from "v0.1"` using the tag specified above. 

Inspect the generated `default.md` file and check that it picked up the commit similar to this:

```
# PAC Changelog

## CS-1

- 045e2a3: CS-1


## Unspecified
```
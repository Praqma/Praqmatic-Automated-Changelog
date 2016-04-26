## Installation instructions for windows

download ruby from [here](http://dl.bintray.com/oneclick/rubyinstaller/rubyinstaller-2.2.4.exe)

during installation let it add ruby to PATH

download [devkit](http://dl.bintray.com/oneclick/rubyinstaller/DevKit-mingw64-32-4.7.2-20130224-1151-sfx.exe)

extract it to `C:\Ruby22\DevKit`

follow these instructions to install it
[https://github.com/oneclick/rubyinstaller/wiki/Development-Kit#quick-start](https://github.com/oneclick/rubyinstaller/wiki/Development-Kit#quick-start)

to install cmake and pkg-config follow the instructions detailed [here](http://stackoverflow.com/a/31254515)

cmake download [URL](http://www.cmake.org/files/v3.2/cmake-3.2.3-win32-x86.zip)

pkg-config download [URL](http://iweb.dl.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win64/Personal%20Builds/ray_linn/64bit-libraries/pkg-config/pkg-config-0.26.7z)

clone PAC from GitHub

run $ gem install bundler

cd to project root directory and run $ bundle


To test if it works

open default_settings.yml

replace task systems with this

```
:task_systems:
  - :name: none
    :regex:
      - { pattern: '/(CS-\d+)/i', label: none }
    :delimiter: '/,|\s/'
```
  
turn off pdf rendering under templates

`- { location: templates/default_html.html, pdf: false, output: default.html }`

set repo location to
`:repo_location: 'C:\projects\Praqmatic-Automated-Changelog\misc'`

create sample repository with one file and one commit. 
The commit message is CS-1.

run $ ruby pac.rb -d {current date}

in the default.md file you should see that it picked up the commit.

PAC id report

Referenced tasks

CS-1
• fb18f56:CS-1

Unreferenced commits


If it does then it worked
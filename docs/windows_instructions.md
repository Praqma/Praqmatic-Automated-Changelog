# Installing PAC on Windows

1. Download and install [Ruby](http://dl.bintray.com/oneclick/rubyinstaller/rubyinstaller-2.2.4.exe), let the installer add Ruby to PATH when given the option.
2. Download [DevKit](http://dl.bintray.com/oneclick/rubyinstaller/DevKit-mingw64-32-4.7.2-20130224-1151-sfx.exe) and install it following [the DevKit quick-start](https://github.com/oneclick/rubyinstaller/wiki/Development-Kit#quick-start).
3. Download and install [cmake](http://www.cmake.org/files/v3.2/cmake-3.2.3-win32-x86.zip) and [pkg-config](http://iweb.dl.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win64/Personal%20Builds/ray_linn/64bit-libraries/pkg-config/pkg-config-0.26.7z) following instructions detailed in [this SO answer](http://stackoverflow.com/a/31254515).
4. Install `bundler` by running `gem install bundler`
5. Clone the pac repository to your local machine: git clone https://github.com/Praqma/Praqmatic-Automated-Changelog.git pac
6. Optionally check-out the latest tag or a specific release tag if you don't want bleeding edge.
7. Change directory to pac (the git clone) and run the command `bundle install` to install all the used Ruby Gems.

## Testing the installation

Open `default_settings.yml` and replace the `:task_systems:` block with the following:

```
:task_systems:
  - :name: none
    :regex:
      - { pattern: '/(CS-\d+)/i', label: none }
    :delimiter: '/,|\s/'
```
  
Disable the pdf generation under the `:templates:` section:

`- { location: templates/default_html.html, pdf: false, output: default.html }`

Set the repository location to the `misc` folder in the PAC project. 

`:repo_location: '...\Praqmatic-Automated-Changelog\misc'`

In the above mentioned folder, create a sample repository containing one commit with the message `CS-1`.

Run `ruby pac.rb -d {current date}` using the date older than your above commit. 

Inspect the generated `default.md` file and check that it picked up the commit similar to this:

```
PAC id report

Referenced tasks

CS-1
• fb18f56:CS-1

Unreferenced commits
```

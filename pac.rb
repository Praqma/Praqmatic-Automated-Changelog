#!/usr/bin/env ruby
# encoding: utf-8
require 'docopt'
require_relative 'lib/core'

doc = <<DOCOPT
Pragmatic changelog 

Usage:
  #{__FILE__} (-d | --date) <to> [<from>] [--outpath=<path>]  [--settings=<settings_file>] [--formats=<format>] [--pattern=<rel_pattern>]     
  #{__FILE__} (-s | --sha) <to> [<from>] [--outpath=<path>]  [--settings=<settings_file>] [--formats=<format>] [--pattern=<rel_pattern>]
  #{__FILE__} (-t | --tag) <to> [<from>] [--outpath=<path>]  [--settings=<settings_file>] [--formats=<format>] [--pattern=<rel_pattern>]  
  #{__FILE__} -h|--help

Options:
  -h --help
             
    Show this screen.
    
  -d --date
             
    Use dates to select the changesets.
     
  -s --sha
              
    Use SHAs to select the changesets.
      
  --output=<path>     
    
    Specify where the changelog should be written.
  
  --settings=<settings_file> 
  
    Path to the settings file used. If nothing is specified default_settings.yml is used
    
  --formats=<format>
  
    Comma delimited string of which formats to use. We currently support html and pdf, markdown is always created
    
  --pattern=<rel_pattern>
  
    Format that describes how your release tags look. This is used together with -t LATEST. We always check agains HEAD/TIP.        
DOCOPT

begin
  require "pp"
  input = Docopt::docopt(doc)
  
  settings_file = File.join(File.dirname(__FILE__), 'default_settings.yml')
  
  unless input['--settings'].nil?
    settings_file = input['--settings']
  end
  
  settings = YAML::load(File.open(settings_file))
  
  unless input['--pattern'].nil? 
    settings[:vcs][:release_regex] = input['--pattern']
  end
  
  Core.load(settings)
  
  if (input['--sha'] || input['-s'])
    changes = Core.vcs.get_commit_messages_by_commit_sha(input['<to>'], input['<from>'])
  elsif (input['--date'] || input['-d'])
    toTime = Core.to_time(input['<to>'])
    unless input['<from>'].nil?
      fromTime = Core.to_time(input['<from>'])
    end    
    changes = Core.vcs.get_commit_messages_by_commit_times(toTime, fromTime)
  else   
    changes = Core.vcs.get_commit_messages_by_tag_name(input['<to>'], input['<from>'])  
  end
  
  #Use core to write the changelog. Depending on the system use we use different ways to do it.
  unless changes == []
    options = []    
    Core.task_system.write_changelog(changes, input['--outpath'])
  else
    puts "No changesets found"
  end 
  #pp Docopt::docopt(doc)
rescue Docopt::Exit => e
  puts e.message
end


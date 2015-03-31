#!/usr/bin/env ruby
#encoding: utf-8
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
    commit_map = Core.vcs.get_commit_messages_by_commit_sha(input['<to>'], input['<from>'])
    changes = commit_map.values
  elsif (input['--date'] || input['-d'])
    toTime = Core.to_time(input['<to>'])
    unless input['<from>'].nil?
      fromTime = Core.to_time(input['<from>'])
    end
    commit_map = Core.vcs.get_commit_messages_by_commit_times(toTime, fromTime)     
    changes = commit_map.values  
  else
    commit_map = Core.vcs.get_commit_messages_by_tag_name(input['<to>'], input['<from>'])    
    changes = commit_map.values
  end

  #Write an ID_REPORT based on task system and regex. This happenes just before we query the actual contents of the task system.   

  Core.task_system.task_id_report(commit_map)

  #Use core to write the changelog. Depending on the system use we use different ways to do it.
  if !Core.task_system.is_a? Task::NoneTaskSystem
    unless changes == []
      options = []    
      Core.task_system.write_changelog(changes, input['--outpath'])
    else
      puts "No changesets found"  
    end 
  end
  #pp Docopt::docopt(doc)
rescue Docopt::Exit => e
  puts e.message
end
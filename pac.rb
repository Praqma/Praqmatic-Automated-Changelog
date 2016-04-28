#!/usr/bin/env ruby
#encoding: utf-8
require 'docopt'
require 'liquid'
require_relative 'lib/core'
require_relative 'lib/report'

doc = <<DOCOPT
Usage:
  #{__FILE__} (-d | --date) <to> [<from>] [--settings=<settings_file>] [--pattern=<rel_pattern>]     
  #{__FILE__} (-s | --sha) <to> [<from>] [--settings=<settings_file>] [--pattern=<rel_pattern>]
  #{__FILE__} (-t | --tag) <to> [<from>] [--settings=<settings_file>] [--pattern=<rel_pattern>]  
  #{__FILE__} -h|--help

Options:
  -h --help
             
    Show this screen.
    
  -d --date
             
    Use dates to select the changesets.
     
  -s --sha
              
    Use SHAs to select the changesets.      
  
  --settings=<settings_file> 
  
    Path to the settings file used. If nothing is specified default_settings.yml is used
    
  --pattern=<rel_pattern>
  
    Format that describes how your release tags look. This is used together with -t LATEST. We always check agains HEAD/TIP         
DOCOPT

begin
  require "pp"

  #Versioning strategy see the file docs/versioning.md
  def version
    if File.exist?(File.dirname(__FILE__)+'/version.stamp')
      version = File.read(File.dirname(__FILE__)+'/version.stamp')
    else
      version = "Unknown version"
    end
    version
  end

  input = Docopt::docopt(doc)
  settings_file = File.join(Dir.pwd, 'default_settings.yml')
  
  unless input['--settings'].nil?
    settings_file = input['--settings']
  end
  
  loaded = YAML::load(File.open(settings_file))

  unless input['--pattern'].nil? 
    loaded[:vcs][:release_regex] = input['--pattern']
  end
  
  #Load the settings
  Core.settings = loaded 

  if (input['--sha'] || input['-s'])
    commit_map = Core.vcs.get_commit_messages_by_commit_sha(input['<to>'], input['<from>'])
  elsif (input['--date'] || input['-d'])
    toTime = Core.to_time(input['<to>'])
    unless input['<from>'].nil?
      fromTime = Core.to_time(input['<from>'])
    end
    commit_map = Core.vcs.get_commit_messages_by_commit_times(toTime, fromTime)     
  else
    commit_map = Core.vcs.get_commit_messages_by_tag_name(input['<to>'], input['<from>'])    
  end

  #This is all our current tasks (PACTaskCollection) Each task is uniquely identified by an ID.
  #We need to iterate each task system
  tasks = Core.task_id_list(commit_map)
  everything_ok = true
  #Apply the task system(s) to each task. Basically populate each task with data from the task system(s)  
  Core.settings[:task_systems].each do |ts|
    everything_ok &= Core.apply_task_system(ts, tasks)
  end

  #Write the ID report (Basically just a list of referenced and non-referenced issues)
  #Takes the list of discovered tasks, and only needs the template settings
  generator = Report::Generator.new
  generator.generate(tasks, commit_map, Core.settings[:templates])
  unless everything_ok
  	if Core.settings[:general][:strict]
  		exit 15
  	else
  	  puts '[PAC] Ignoring encountered errors. Strict mode is disabled.'
  	  exit 0
  	end
  end

rescue Docopt::Exit => e
  puts "Praqmatic Automated Changelog (PAC)"
  puts "#{version}\n"
  puts e.message
end

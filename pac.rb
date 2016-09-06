#!/usr/bin/env ruby
#encoding: utf-8
require 'docopt'
require 'liquid'
require_relative 'lib/core'
require_relative 'lib/report'

doc = <<DOCOPT
Usage:
  #{__FILE__} (-d | --date) <to> [<from>] [options] [-v...] [-q...]
  #{__FILE__} (-s | --sha) <to> [<from>] [options] [-v...] [-q...]
  #{__FILE__} (-t | --tag) <to> [<from>] [options] [-v...] [-q...]
  #{__FILE__} from <oldest-ref> [--to <newest-ref> | --to-date <newest-date>] [options] [-v...] [-q...] [-c (<user> <password> <target>)]... 
  #{__FILE__} from-date <from_date> [--to <newest-ref> | --to-date <newest-date>] [options] [-v...] [-q...] [-c (<user> <password> <target>)]...
  #{__FILE__} from-latest-tag <approximation> [--to <newest-ref> | --to-date <newest-date>] [options] [-v...] [-q...] [-c (<user> <password> <target>)]...
  #{__FILE__} -h|--help

Options:
  -h --help  Show this screen.

  --from <oldest-ref>  Specify where to stop searching for commit. For git this takes anything that rev-parse accepts. Such as HEAD~3 / Git sha or tag name.

  --from-date <from-date>  Use dates to select the changesets.

  --from-latest-tag  Looks for the newest commit that the tag with <approximation> points to.  
    
  -d --date  Use dates to select the changesets.  

  -s --sha  Deprecated since 2.1.0. Use --from instead.
              
  --settings=<path>  Path to the settings file used. If nothing is specified default_settings.yml is used      

  --properties=<properties>  

    Allows you to pass in additional variables to the Liquid templates. Must be in JSON format. Namespaced under properties.* in 
    your Liquid templates. Referenced like so '{{properties.[your-variable]}}' in your templates.

    JSON keys and values should be wrapped in quotation marks '"' like so: --properties='{ "title":"PAC Changelog" }'      

  -v  More verbose output. Can be repeated to increase output verbosity or to cancel out -q

  -q  Less verbose output. Can be repeated for more silence or to cancel out -v
DOCOPT

begin
  require "pp"

  #Versioning strategy see the file docs/versioning.md
  #We need to read where the symlink points to because the docker container links pac to /usr/bin
  def version
    path = File.symlink?(__FILE__) ? File.readlink(__FILE__) : __FILE__
    dir = File.dirname(path)
    if File.exist?(dir+'/version.stamp')
      version = File.read(dir+'/version.stamp')
    else
      version = "Unknown version"
    end
    version
  end

  input = Docopt::docopt(doc)

  #This should print out any and all errors related to creating settings for PAC. This captures
  #JSON parser errors.
  begin
    configuration = Core.read_settings_file(input)
    loaded = Core.generate_settings(input, configuration)
    #Load the settings
    Core.settings = loaded 
  rescue JSON::ParserError => pe
    Logging.verboseprint(0, "[PAC] Error paring JSON from --properties switch")
    Logging.verboseprint(0, "[PAC] Exception caught while parsing command line options: #{pe}")
    exit 6
  rescue Exception => e
    Logging.verboseprint(0, "[PAC] Unspecified error caught while creating configuration")
    Logging.verboseprint(0, "[PAC] Exception caught while creating configuration: #{e}")
    exit 7
  end

  if (input['from'])
    if input['<newest-date>']
      Logging.verboseprint(0,"[PAC] Fiding commits until date #{input['<newest-date>']}")
      to_commit = Core.vcs.get_first_commit_after(Core.to_time(input['<newest-date>']))
      Logging.verboseprint(0, "[PAC] Finding commits up and until commit: #{to_commit}")
      commit_map =   Core.vcs.get_delta(input['<oldest-ref>'], to_commit)
    else
      commit_map = Core.vcs.get_delta(input['<oldest-ref>'], input['<newest-ref>'])
    end   
  elsif (input['from-latest-tag'])
    found_tag = Core.vcs.get_latest_tag(input['<approximation>'])
    commit_map = Core.vcs.get_delta(found_tag, input['<newest-ref>'])
  elsif input['from-date']
    if input['--to-date']
      Logging.verboseprint(0, "[PAC] Using dates. Finding commits between #{input['<newest-date>']} and #{input['<from-date>']}")
    elsif input['--to']
      Logging.verboseprint(0, "[PAC] Using dates. Finding commits between commit #{input['newest-ref']} and #{input['<from-date>']}")      
    else
      Logging.verboseprint(0, "[PAC] Using dates. Finding commits from tip of repository until #{input['<from-date>']}")
    end
  elsif (input['--sha'] || input['-s'])
    Logging.verboseprint(0, "[PAC] Warning: Using deprecated method call. Use --from instead")
    commit_map = Core.vcs.get_delta(input['<to>'], input['<from>'])
  elsif (input['--date'] || input['-d'])
    toTime = Core.to_time(input['<to>'])
    unless input['<from>'].nil?
      fromTime = Core.to_time(input['<from>'])
    end
    commit_map = Core.vcs.get_commit_messages_by_commit_times(toTime, fromTime)     
  else
    Logging.verboseprint(0, "[PAC] Warning: Using deprecated method call. Use --from instead. This method accepts both tags and git shas")
    commit_map = Core.vcs.get_delta(input['<to>'], input['<from>'])    
  end

  #This is all our current tasks (PACTaskCollection) Each task is uniquely identified by an ID.
  #We need to iterate each task system
  tasks = Core.task_id_list(commit_map)
  everything_ok = true
  #Apply the task system(s) to each task. Basically populate each task with data from the task system(s)  
  Core.settings[:task_systems].each do |ts|
    res = Core.apply_task_system(ts, tasks)
    everything_ok &= res 
  end

  #Write the ID report (Basically just a list of referenced and non-referenced issues)
  #Takes the list of discovered tasks, and only needs the template settings
  generator = Report::Generator.new(tasks, commit_map)
  generator.generate(Core.settings)
  unless everything_ok
  	if Core.settings[:general][:strict]
  		exit 15
  	else
  	  Logging.verboseprint(1, '[PAC] Ignoring encountered errors. Strict mode is disabled.')
  	  exit 0
  	end
  end

rescue Docopt::Exit => e
  puts "Praqmatic Automated Changelog (PAC)"
  puts "#{version}\n"
  puts e.message
  puts
end

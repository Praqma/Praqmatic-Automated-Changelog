# encoding: utf-8
require 'yaml'
require_relative "./task.rb"
require_relative "./gitvcs"
require_relative "./logging"

module Core extend self

  def cli_text(file)
    cli = <<DOCOPT
    Usage:
      pac from <oldest-ref> [--to=<newest-ref>] [options] [-v...] [-q...]
      pac from-latest-tag <approximation> [--to=<newest-ref>] [options] [-v...] [-q...]
      pac -h|--help

    Options:
      -h --help  Show this screen.

      --no-ssl-verify  Disables peer verification over https

      --ssl-verify  Enables peer verification over https

      --strict  Enables strict mode. Will exit nonzero on error fetching data

      --no-strict  Disables strict mode. Will ignore network errors while getting issues (Common issue: Misspelled issue number)

      --from <oldest-ref>  Specify where to stop searching for commit. For git this takes anything that rev-parse accepts. Such as HEAD~3 / Git sha or tag name.

      --from-latest-tag  Looks for the newest commit that the tag with <approximation> points to.

      --settings=<path>  Path to the settings file used. If nothing is specified default_settings.yml is used

      --properties=<properties>

        Allows you to pass in additional variables to the Liquid templates. Must be in JSON format. Namespaced under properties.* in
        your Liquid templates. Referenced like so '{{properties.[your-variable]}}' in your templates.

        JSON keys and values should be wrapped in quotation marks '"' like so: --properties='{ "title":"PAC Changelog" }'

      -v  More verbose output. Can be repeated to increase output verbosity or to cancel out -q

      -q  Less verbose output. Can be repeated for more silence or to cancel out -v
DOCOPT
    cli
  end

  def settings
    if defined?(@@settings).nil?
      {}
    else
      @@settings
    end
  end

  def settings=(val)
    @@settings = val
  end

  #Reads the command line options. And based on this it will return the
  #path of the settings file to use.
  def read_settings_file(input)
    settings_file = File.join(Dir.pwd, 'settings/default_settings.yml')
    unless input['--settings'].nil?
      settings_file = input['--settings']
    end

    unless File.exists?(settings_file)
      raise "Settings file '#{settings_file}' does not exist"
    end

    File.read(settings_file)
  end

  #Creates the final settings based on additonal command line arguments
  #Parameters
  # cmdline       - The command line arguments parsed by docopt. Essentially a ruby hash
  # configuration - The contents of the settings file
  #Exceptions
  # All exections are throw as is. Handled in pac.rb.
  def generate_settings(cmdline, configuration)
    loaded = YAML::load(configuration)

    if loaded[:properties].nil?
      loaded[:properties] = {}
    end

    if loaded[:general].nil?
      loaded[:general] = { }
    end

    if loaded[:general][:ssl_verify].nil?
      loaded[:general][:ssl_verify] = true
    end

    if cmdline['--strict']
      loaded[:general][:strict] = true
    end

    if cmdline['--no-strict']
      loaded[:general][:strict] = false
    end

    if cmdline['--ssl-verify']
      loaded[:general][:ssl_verify] = true
    end

    if cmdline['--no-ssl-verify']
      loaded[:general][:ssl_verify] = false
    end

    unless cmdline['--properties'].nil?
      json_value = JSON.parse(cmdline['--properties'])
      loaded[:properties] = loaded[:properties].merge(json_value)
    end
    loaded[:verbosity] = Logging.calc_verbosity(cmdline)
    loaded
  end

  #Requires a configuration section for the task system to be applied
  def apply_task_system(task_system, tasks)
    val = true
    Logging.verboseprint(1, "[PAC] Applying task system #{task_system[:name]}")
    if task_system[:name] == 'json'
      val = Task::JsonTaskSystem.new(task_system).apply(tasks)
    end
    val
  end

  def vcs
    Vcs::GitVcs.new(settings[:vcs])
  end

  #This is now core functionality. The task of generating a collection of tasks based on the commits found
  #This takes in a PACCommitCollection and returns a PACTaskCollection
  def task_id_list(commits)
    regex_arr = []
    tasks = Model::PACTaskCollection.new
    commits.each do |c_pac|
      referenced = false
      #Regex ~ Eacb regex in the task system
      settings[:task_systems].each do |ts|
        #Loop over each task system. Parse commits for matches
        if ts.has_key? :delimiter
          split_pattern = eval(ts[:delimiter])
        end
        if ts.has_key? :regex
          tasks_for_commit = c_pac.matchtask(ts[:regex], split_pattern)
          tasks_for_commit.each do |t|
            t.applies_to = ts[:name]
          end
          #If a task was found
          unless tasks_for_commit.empty?
            referenced = true
            tasks.add(tasks_for_commit)
          end
        end
      end

      if !referenced
        task = Model::PACTask.new
        task.add_commit(c_pac)
        tasks.add(task)
      end
    end

    tasks
  end
end

# encoding: utf-8
require 'yaml'
require_relative "./task.rb"
require_relative "./vcs"
require_relative "./gitvcs"
require_relative "./mercurialvcs"

module Core extend self
   
  attr_accessor :settings

  def load(file)
    if file.is_a? String     
      @settings = YAML::load(File.open(file))
    elsif file.is_a? Hash
      @settings = file
    else
      raise ArgumentError, 'The settings passed into this must be either a string or an array of arguments'
    end    
  end
  
  def task_system
    if @settings.include? :fogbugz
      Task::FogBugzTaskSystem.new(@settings)
    elsif @settings.include? :trac
      Task::TracTaskSystem.new(@settings)
    else 
      Task::NoneTaskSystem.new(@settings)
    end
  end
  
  def vcs
    if @settings[:vcs]['type'] == 'git'
      Vcs::GitVcs.new(@settings[:vcs])
    elsif @settings[:vcs]['type'] == 'hg'
      Vcs::MercurialVcs.new(@settings[:vcs])
    else
      raise ArgumentError, 'The configuration settings does not include any supported (d)vcs'
    end
  end
  
  def to_time(datestring)
    DateTime.strptime(datestring, @settings[:general]['date_template']).to_time    
  end

end

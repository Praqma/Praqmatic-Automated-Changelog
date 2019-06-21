# encoding: utf-8
require_relative "./core.rb"

module Logging extend self
  def calc_verbosity(input)
    loudness = 0
    v = input['-v']
    q = input['-q']
    if v then loudness = loudness + v end
    if q then loudness = loudness - q end
    loudness
  end

  def verboseprint(fromlevel, str)
    str = v(fromlevel, str)
    puts str unless str.nil?
  end

  def v(fromlevel, str)
    if Core.settings[:verbosity].nil?
      return str
    end
    if Core.settings[:verbosity] >= fromlevel
      return str
    end
  end
end

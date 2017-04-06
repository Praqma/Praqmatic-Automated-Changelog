require 'test/unit'
require_relative "../../lib/model.rb"
require_relative "../../lib/task.rb"
require_relative "../../lib/decorators.rb"
require_relative "../../lib/core.rb"
require_relative "../../lib/report.rb"

class GeneratorTest < Test::Unit::TestCase

	#Setup everything needed to generate the template
  Test::Unit.at_start do
  	@@generator = Report::Generator.new

  	@@c_noprops = <<-eos
:general:
 
:templates:
  - { location: templates/default_id_report.md }

:task_systems:
  - 
    :name: none
    :regex:
      - { pattern: '/Issue:\s*(\d+)/i', label: none }
:vcs:
  :type: git
  :repo_location: '.'
eos

  	@@c_props = <<-eos
:general:
:properties:
  title: 'Awesome Changelog Volume 42'
  extra: 'Extra Stuff' 
:templates:
  - { location: templates/default_id_report.md}
:task_systems:
  - 
    :name: none
    :regex:
      - { pattern: '/Issue:\s*(\d+)/i', label: none }
:vcs:
  :type: git
  :repo_location: '.'
eos

  	@@min_template = '# {{properties.title}}'

  end

	#Test that show properties from command line are exposed and printed in the template
	def test_properties_from_cmd
		p = { '--properties' =>  '{"title" : "Awesome Changelog Volume 2" }' }		
		#Generate the settings, using command line arguments p and the config file defined in c_noprops
		settings = Core.generate_settings(p, @@c_noprops)

		#Extract the properties
		properties = @@generator.to_liquid_properties(settings)

		#Render the template. At this point we just need to make sure that the property properties.title has been substituted
		render = @@generator.render_template(@@min_template, properties)

		expected = '# Awesome Changelog Volume 2'
		assert_equal(expected, render)
	end

	#Test that verified that configuration file properties are exposed and printed in the template. 
	def test_properties_from_config		
		p = {}	
		#Generate the settings, using command line arguments p and the config file defined in c_props
		settings = Core.generate_settings(p, @@c_props)

		#Extract the properties
		properties = @@generator.to_liquid_properties(settings)

		#Render the template. At this point we just need to make sure that the property properties.title has been substituted
		render = @@generator.render_template(@@min_template, properties)

		expected = '# Awesome Changelog Volume 42'
		assert_equal(expected, render)		
	end

	#We need to make sure that parameters specified via the command line overrides the values specified in the settings file. 	
	def test_properties_precedence
		p = { '--properties' =>  '{"title" : "Awesome Changelog Volume 2" }'}
		
		#Generate the settings, using command line arguments p and the config file defined in c_props
		settings = Core.generate_settings(p, @@c_props)

		#Extract the properties
		properties = @@generator.to_liquid_properties(settings)

		#Render the template. At this point we just need to make sure that the property properties.title has been substituted
		render = @@generator.render_template(@@min_template, properties)

		expected = '# Awesome Changelog Volume 2'
		assert_equal(expected, render)				
	end

	#We also need to make sure that partial overrides work. That is, if only one of the values specified in the configuration 
	#has been overridden. In the case below we still want to keep the 'extra' property, but we need to override the 'title'.  
	def test_partial_override
		p = { '--properties' =>  '{"title" : "Awesome Changelog Volume 43" }'}
		
		#Generate the settings, using command line arguments p and the config file defined in c_props
		settings = Core.generate_settings(p, @@c_props)

		#Extract the properties
		properties = @@generator.to_liquid_properties(settings)

		#Render the template. At this point we just need to make sure that the property properties.title has been substituted
		extra_template = '{{properties.title}} and {{properties.extra}}'
		render = @@generator.render_template(extra_template, properties)

		expected = 'Awesome Changelog Volume 43 and Extra Stuff'
		assert_equal(expected, render)		
	end

end
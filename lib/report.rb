module Report
	class Generator
		require 'liquid'
		#Constructor for the generator. 
		#Arguments:
	  # tasks     - A 'PACTaskCollection' of populated tasks
	  # commits   - A 'PACCommitCollection' list of commits. used to tally those that were referenced and those that were not
		def initialize(tasks = Model::PACTaskCollection.new, commits = Model::PACCommitCollection.new)
			@tasks = tasks
			@commits = commits
		end 

	  #Generates the output files specifed in the config file.
	  #Arguments:
	  # config 		- The configuration used 
	  def generate(config)
	    config[:templates].each do |t|
	    	unless t['output'].nil?        
		      File.open(t['output'],'w:UTF-8') do |file|
		      	file << render_template(File.read(t['location']), to_liquid_properties(config)) 
		        File.chmod(0777, file)
		      end

		      if t['pdf'] == true
		        output_file_path = t['output'].sub(/\.html$/, '.pdf')
		        kit = PDFKit.new(File.new(t['output']), :page_size => 'A4')
		        kit.to_file(output_file_path)
		        File.chmod(0777, output_file_path)
		      end
		    else
		    	puts "========== #{t['location']} =========="
		    	puts render_template(File.read(t['location']), to_liquid_properties(config))
		    	footer = "=" * (t['location'].length + 22)
		    	puts "#{footer}"
	    	end
	    end
	  end

	  #Argments:
	  # template - A string representing the template to render
	  # props    - A ruby has with additional properties
	  #Returns - The rendered template	  
	  def render_template(template, props = {})
			Liquid::Template.parse(template).render( { 
        'tasks' => @tasks, 	     
        'pac_c_count' => @commits.count,
        'pac_c_referenced' => @commits.count_with,
        'pac_health' => @commits.health,
        'pac_c_unreferenced' => @commits.count_without  
      }.merge!(props))
	  end

	  #Convert the properties from the parsed yml into string keys so that Liquid can read it. 
	  def to_liquid_properties(config)
	  	std = { 'properties' => { } }
	  	unless config[:properties].nil?
	  		std['properties'].merge!(config[:properties])
	  	end
	  	std
	  end 
	end
end
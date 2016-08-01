# -*- coding: utf-8; -*-
module PAC__TestCases_Tags
  require 'pp'
  require 'fileutils'
  require 'zip/zip'
  require 'ci/reporter/rake/test_unit_loader.rb'

  class GitTagging < Test::Unit::TestCase
    require_relative '../../../lib/core'
    require 'fileutils'
    require 'open3'
    # Helper function to unzip our repo
    def unzip_file (file, destination)
      Zip::ZipFile.open(file) { |zip_file|
        zip_file.each { |f|
          f_path=File.join(destination, f.name)
          FileUtils.mkdir_p(File.dirname(f_path))
          zip_file.extract(f, f_path) unless File.exist?(f_path)
        }
        puts "Successfully extracted"
      }
    end

    # Test setup: executed before each test
    def setup
      unzip_file("test/resources/repo_with_tags.zip", "test/resources")
    end

    # Executed after each test
    def teardown
      FileUtils.rm_rf("test/resources/repo_with_tags")
    end

    def test_get_latest_tag_function
      require_relative '../../../lib/core'
      settings_file = File.join(File.dirname(__FILE__), '../../resources/repo_with_tags_settings.yml')
      settings = YAML::load(File.open(settings_file)) 
      settings[:verbose] = 1
      Core.settings = settings
      tag_name = Core.vcs.get_latest_tag(nil)     
      #this should be 'second'
      assert_equal("second", tag_name)

      tag_first = Core.vcs.get_latest_tag "first"
      assert_equal("first", tag_first)
    end

  end # class
end # module


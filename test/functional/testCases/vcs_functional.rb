module PAC__TestCases_IncludesSince
  require 'fileutils'
  require 'zip/zip'
  gem 'test-unit'
  require 'test/unit'

  # Issue documented in FogBugz case 13433
  #
  # When generating a changelog since a certain tag/sha, the commit of that tag/sha was included.
  # This test asserts that the commit of the 'since' (AKA tail) tag/sha isn't included in the changelog.

  class IncludesSince < Test::Unit::TestCase

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
      unzip_file("test/resources/IncludesSinceTestRepository.zip", "test/resources")
    end

    # Executed after each test
    def teardown
      FileUtils.rm_rf("test/resources/IncludesSinceTestRepository")
    end

    def test_includesSince_tag
      require_relative '../../../lib/core'
      settings_file = File.join(File.dirname(__FILE__), '../../resources/IncludesSinceTestRepository_settings.yml')

      settings = YAML::load(File.open(settings_file))

      Core.settings = settings
      #Returns a PACCommitCollection
      commit_map = Core.vcs.get_delta('v1.0') # Since Commit 1
      # These are the commits AFTER the 1.0 tag
      assert_true(commit_map.commits.include?(Model::PACCommit.new("dfd8fd3ade868072f31701aac6fba1f2cf965dd7"))) # Commit 3
      assert_true(commit_map.commits.include?(Model::PACCommit.new("62f42e26a401524637839ba4ff969194303cff7c"))) # Commit 2
      # This is the commit OF the 1.0Modelg
      assert_true(!commit_map.commits.include?(Model::PACCommit.new("7164ed9cec32195e44c7f2a2abd764df37921863"))) # Commit 1
    end

    def test_includesSince_sha
      require_relative '../../../lib/core'
      settings_file = File.join(File.dirname(__FILE__), '../../resources/IncludesSinceTestRepository_settings.yml')

      settings = YAML::load(File.open(settings_file))

      Core.settings = settings

      commit_map = Core.vcs.get_delta('7164ed9cec32195e44c7f2a2abd764df37921863') # Since Commit 1
      # These are the commits AFTER the 1.0 tag
      assert_true(commit_map.commits.include?(Model::PACCommit.new("dfd8fd3ade868072f31701aac6fba1f2cf965dd7"))) # Commit 3
      assert_true(commit_map.commits.include?(Model::PACCommit.new("62f42e26a401524637839ba4ff969194303cff7c"))) # Commit 2
      # This is the commit OF the 1.0 tag
      assert_true(!commit_map.commits.include?(Model::PACCommit.new("7164ed9cec32195e44c7f2a2abd764df37921863"))) # Commit 1
    end

  end # class
end # module
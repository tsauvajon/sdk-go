require 'kuzzle'

Kuzzle.new "localhost"
puts "#{Kuzzle.server.now()}"

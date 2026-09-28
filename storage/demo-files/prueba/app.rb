class Greeter
  def initialize(name)
    @name = name
  end

  def message
    "Hola, #{@name}"
  end
end

%w[Ruby Go PHP].each do |language|
  puts Greeter.new(language).message
end

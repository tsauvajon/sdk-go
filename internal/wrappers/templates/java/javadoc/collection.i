%typemap(javaimports) kuzzleio::Collection "
/* The type Collection. */"

%javamethodmodifiers kuzzleio::Collection::Collection(Kuzzle *kuzzle) "
  /**
   * Constructor
   *
   * @param kuzzle  Kuzzle instance
   */
  public";

%javamethodmodifiers kuzzleio::Collection::create(const std::string& index, const std::string &collection, query_options* options) "
  /**
   * Create a new empty data collection, with no associated mapping.
   * Kuzzle automatically creates data collections when storing documents, but there are cases where we want to create and prepare data collections before storing documents in it.
   *
   * @param index - Index where to create the collection
   * @param collection - The name of the collection
   * @param options - Request options
   */
  public";

%javamethodmodifiers kuzzleio::Collection::create(const std::string& index, const std::string &collection) "
  /**
   * {@link #create(String index, String collection, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::exists(const std::string& index, const std::string &collection, query_options* options) "
  /**
   * Check if a collection exists
   *
   * @param index - Index where to create the collection
   * @param collection - The name of the collection
   * @param options - Request options
   * @return a boolean
   */
  public";

%javamethodmodifiers kuzzleio::Collection::exists(const std::string& index, const std::string &collection) "
  /**
   * {@link #exists(String index, String collection, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::list(const std::string& index, query_options* options) "
  /**
   * List collections
   *
   * @param index - Parent data index name
   * @param options - Request options
   * @return a json containing the list of collections
   */
  public";

%javamethodmodifiers kuzzleio::Collection::list(const std::string& index) "
  /**
   * {@link #list(String index, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::truncate(const std::string& index, const std::string &collection, query_options* options) "
  /**
   * Truncate the data collection, removing all stored documents but keeping all associated mappings.
   *
   * @param index - Parent data index name
   * @param collection - The name of the collection
   * @param options - Request options
   */
  public";

%javamethodmodifiers kuzzleio::Collection::truncate(const std::string& index, const std::string &collection) "
  /**
   * {@link #truncate(String index, String collection, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::getMapping(const std::string& index, const std::string &collection, query_options* options) "
  /**
   * Get the mapping for this collection
   *
   * @param index - Parent data index name
   * @param collection - The name of the collection
   * @param options - Request options
   */
  public";

%javamethodmodifiers kuzzleio::Collection::getMapping(const std::string& index, const std::string &collection) "
  /**
   * {@link #getMapping(String index, String collection, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::updateMapping(const std::string& index, const std::string &collection, const std::string &body, query_options* options) "
  /**
   * Update the mapping for this collection
   *
   * @param index - Parent data index name
   * @param collection - The name of the collection
   * @param body - The json representing the mapping
   * @param options - Request options
   */
  public";

%javamethodmodifiers kuzzleio::Collection::updateMapping(const std::string& index, const std::string &collection, const std::string &body) "
  /**
   * {@link #updateMapping(String index, String collection, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::getSpecifications(const std::string& index, const std::string &collection, query_options* options) "
  /**
   * Get the specifications for this collection
   *
   * @param index - Parent data index name
   * @param collection - The name of the collection
   * @param options - Request options
   * @return specifications
   */
  public";

%javamethodmodifiers kuzzleio::Collection::getSpecifications(const std::string& index, const std::string &collection) "
  /**
   * {@link #getSpecifications(String index, String collection, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::searchSpecifications(const std::string &filters, query_options* options) "
  /**
   * Searches specifications across indexes/collections according to the provided filters.
   *
   * @param filters - The json representing the filters
   * @param options - Request options
   * @return a SearchResult
   */
  public";

%javamethodmodifiers kuzzleio::Collection::searchSpecifications(const std::string &filters) "
  /**
   * {@link #searchSpecifications(String filters, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::updateSpecifications(const std::string &index, const std::string &collection, const std::string &body, query_options* options) "
  /**
   * Updates the current specifications of the specified collection
   *
   * @param index - Parent data index name
   * @param collection - collection to update
   * @param body - The new specifications
   * @param options - Request options
   */
  public";

%javamethodmodifiers kuzzleio::Collection::updateSpecifications(const std::string &index, const std::string &collection, const std::string &body) "
  /**
   * {@link #updateSpecifications(String index, String collection, String body, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::validateSpecifications(const std::string &body, query_options* options) "
  /**
   * Validates the provided specifications
   *
   * @param body - The json representing the specifications
   * @param options - Request options
   */
  public";

%javamethodmodifiers kuzzleio::Collection::validateSpecifications(const std::string &body) "
  /**
   * {@link #validateSpecifications(String body, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Collection::deleteSpecifications(const std::string& index, const std::string &collection, query_options* options) "
  /**
   * Delete the specifications for this collection
   *
   * @param index - Parent data index name
   * @param collection - The name of the collection
   * @param options - Request options
   */
  public";

%javamethodmodifiers kuzzleio::Collection::deleteSpecifications(const std::string& index, const std::string &collection) "
  /**
   * {@link #deleteSpecifications(String index, String collection, QueryOptions options)}
   */
  public";
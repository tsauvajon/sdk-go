%typemap(javaimports) kuzzleio::Server "
/* The type Server. */"

%javamethodmodifiers kuzzleio::Server::adminExists(query_options* options) "
  /**
   * Check if an admin exists in kuzzle
   *
   * @param options - Request options
   * @return a boolean
   */
  public";

%javamethodmodifiers kuzzleio::Server::adminExists() "
  /**
   * {@link #adminExists(QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Server::getAllStats(query_options*) "
  /**
   * Get all Kuzzle usage statistics frames
   *
   * @param options - Request options
   * @return a AllStatisticsResult
   */
  public";

%javamethodmodifiers kuzzleio::Server::getAllStats() "
  /**
   * {@link #getAllStatistics(QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Server::getStats(unsigned long start, unsigned long end, query_options *options) "
  /**
   * Get Kuzzle usage statistics
   *
   * @param start - Starting timestamp
   * @param end - Ending timestamp
   * @param options - Request options
   * @return a json representing statistics
   */
  public";

%javamethodmodifiers kuzzleio::Server::getStats(unsigned long start, unsigned long end) "
  /**
   * {@link #getStats(long start, long end, QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Server::getLastStats(query_options *options) "
  /**
   * Get last Kuzzle usage statistics.
   * By default, snapshots are made every 10 seconds and they are stored for 1 hour.
   *
   * @param options - Request options
   * @return a json representing last statistics
   */
  public";

%javamethodmodifiers kuzzleio::Server::getLastStats() "
  /**
   * {@link #getLastStats(QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Server::getConfig(query_options *options) "
  /**
   * Get the current Kuzzle configuration.
   * By default, snapshots are made every 10 seconds and they are stored for 1 hour.
   *
   * @param options - Request options
   * @return a json representing last statistics
   */
  public";

%javamethodmodifiers kuzzleio::Server::getConfig() "
  /**
   * {@link #getConfig(QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Server::info(query_options*) "
  /**
   * Gets information about Kuzzle plugins and active services.
   *
   * @param options - Request options
   * #return a json representing informations
   */
  public";

%javamethodmodifiers kuzzleio::Server::info() "
  /**
   * {@link #info(QueryOptions)}
   */
  public";
%typemap(javaimports) kuzzleio::Document "
/* The type Document. */"

%javamethodmodifiers kuzzleio::Document::count(const std::string& index, const std::string& collection, const std::string& body, query_options *options) "
  /**
   * Count the number of document in a given collection according to a filter
   *
   * @param index - the parent index
   * @param collection - the collection name
   * @param filters - the filters
   * @param body - the content of the document
   * @param options - Request options
   * @return number of document
   */
  public";

%javamethodmodifiers kuzzleio::Document::count(const std::string& index, const std::string& collection, const std::string& body) "
  /**
   * {@link #count(String index, String collection, String filters, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::create(const std::string& index, const std::string& collection, const std::string& id, const std::string& body, query_options *options) "
  /**
   * Create a new document in kuzzle
   *
   * @param index - the index where to create the document
   * @param collection - the collection where to create the document
   * @param id - the document id
   * @param body - the content of the document
   * @param options - Request options
   * @return document id
   */
  public";

%javamethodmodifiers kuzzleio::Document::create(const std::string& index, const std::string& collection, const std::string& id, const std::string& body) "
  /**
   * {@link #create(String index, String collection, String id, String body, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::createOrReplace(const std::string& index, const std::string& collection, const std::string& id, const std::string& body, query_options *options) "
  /**
   * Creates a new document in the persistent data storage, or replace it if it already exists.
   *
   * @param index - the index where to create the document
   * @param collection - the collection where to create the document
   * @param id - the document id
   * @param body - the content of the document
   * @param options - Request options
   * @return document id
   */
  public";

%javamethodmodifiers kuzzleio::Document::createOrReplace(const std::string& index, const std::string& collection, const std::string& id, const std::string& body) "
  /**
   * {@link #createOrReplace(String index, String collection, String id, String body, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::delete_(const std::string& index, const std::string& collection, const std::string& id, query_options *) "
  /**
   * Delete this document from Kuzzle
   *
   * @param index - the index where to delete the document
   * @param collection - the collection where to delete the document
   * @param id - the document id
   * @param options - Request options
   * @return string id of the deleted document
   */
  public";

%javamethodmodifiers kuzzleio::Document::delete_(const std::string& index, const std::string& collection, const std::string& id) "
  /**
   * {@link #delete_(String, String, String, QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::exists(const std::string& index, const std::string& collection, const std::string& id, query_options *) "
  /**
   * Ask Kuzzle if this document exists
   * 
   * @param index - parent index
   * @param collection - the collection
   * @param id - Id of the document
   * @param options - Request options
   * @return bool
   */
  public";

%javamethodmodifiers kuzzleio::Document::exists(const std::string& index, const std::string& collection, const std::string& id) "
  /**
   * {@link #exists(String, String, QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::deleteByQuery(const std::string& index, const std::string& collection, const std::string& body, query_options *) "
  /**
   * Deletes all the documents from Kuzzle that match the given filter or query.
   *
   * @param index - the index where to delete the document
   * @param collection - the collection where to delete the document
   * @param body - A set of filters or queries matching documents you are looking for
   * @param options - Request options
   * @return array if id's
   */
  public";

%javamethodmodifiers kuzzleio::Document::deleteByQuery(const std::string& index, const std::string& collection, const std::string& body) "
  /**
   * {@link #deleteByQuery(String, String, String, QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::get(const std::string& index, const std::string& collection, const std::string& id, query_options *) "
  /**
   * Retrieves the corresponding document from the database.
   *
   * @param index - the index from where to get the document
   * @param collection - the collection from where to get the document
   * @param id - the document id
   * @param options - Request options
   * @return json representing the document
   */
  public";

%javamethodmodifiers kuzzleio::Document::get(const std::string& index, const std::string& collection, const std::string& id) "
  /**
   * {@link #get(String, String, String, QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::replace(const std::string& index, const std::string& collection, const std::string& id, const std::string& body, query_options *options) "
  /**
   * Replaces an existing document in the persistent data storage
   *
   * @param index - the index where to replace the document
   * @param collection - the collection where to replace the document
   * @param id - the document id
   * @param body - the content of the document
   * @param options - Request options
   * @return document id
   */
  public";

%javamethodmodifiers kuzzleio::Document::replace(const std::string& index, const std::string& collection, const std::string& id, const std::string& body) "
  /**
   * {@link #createOrReplace(String index, String collection, String id, String body, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::update(const std::string& index, const std::string& collection, const std::string& id, const std::string& body, query_options* options) "
  /**
   * Update parts of a document
   *
   * @param index - the index where to update the document
   * @param collection - the collection where to update the document
   * @param id - Document unique identifier
   * @param content - Document content to update
   * @param options - Request options
   * @return the document
   */
  public";

%javamethodmodifiers kuzzleio::Document::update(const std::string& index, const std::string& collection, const std::string& id, const std::string& body) "
  /**
   * {@link #update(String index, String collection, String id, String body)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::validate(const std::string& index, const std::string& collection, const std::string& body, query_options* options) "
  /**
   * Validates data against existing validation rules
   *
   * @param index - the index where to validate the document
   * @param collection - the collection where to validate the document
   * @param body - Document content to update
   * @param options - Request options
   * @return the document
   */
  public";

%javamethodmodifiers kuzzleio::Document::validate(const std::string& index, const std::string& collection, const std::string& body) "
  /**
   * {@link #validate(String index, String collection, String body)}
   */
  public";
  
%javamethodmodifiers kuzzleio::Document::search(const std::string& index, const std::string& collection, const std::string& body, query_options* options) "
  /**
   * Validates data against existing validation rules
   *
   * @param index - the index where to search
   * @param collection - the collection where to search
   * @param body - the filters to apply on the search
   * @param options - Request options
   * @return a SearchResult object
   */
  public";

%javamethodmodifiers kuzzleio::Document::search(const std::string& index, const std::string& collection, const std::string& body) "
  /**
   * {@link #search(String index, String collection, String body)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::mCreate(const std::string& index, const std::string& collection, const std::string& body, query_options *options) "
  /**
   * Creates new documents in the persistent data storage.
   *
   * @param index - the index where to create the documents
   * @param collection - the collection where to create the documents
   * @param body - json representing the documents
   * @param options - Request options
   * @return documents
   */
  public";

%javamethodmodifiers kuzzleio::Document::mCreate(const std::string& index, const std::string& collection, const std::string& body) "
  /**
   * {@link #mCreate(String index, String collection,  String body, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::mCreateOrReplace(const std::string& index, const std::string& collection, const std::string& body, query_options *options) "
  /**
   * Creates or replaces documents in the persistent data storage.
   *
   * @param index - the index where to create or replace the documents
   * @param collection - the collection where to create or replace the document
   * @param body - the content of the document
   * @param options - Request options
   * @return documents
   */
  public";

%javamethodmodifiers kuzzleio::Document::mCreateOrReplace(const std::string& index, const std::string& collection, const std::string& body) "
  /**
   * {@link #mCreateOrReplace(String index, String collection, String id, String body, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::mDelete(const std::string& index, const std::string& collection, const std::vector<std::string>& ids, query_options *) "
  /**
   * Delete documents from Kuzzle
   *
   * @param index - the index where to delete the documents
   * @param collection - the collection where to delete the documents
   * @param ids - the documents id
   * @param options - Request options
   * @return a StringVector object
   */
  public";

%javamethodmodifiers kuzzleio::Document::mDelete(const std::string& index, const std::string& collection, const std::vector<std::string>& ids) "
  /**
   * {@link #mDelete(String, String, StringVector, QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::mGet(const std::string& index, const std::string& collection, const std::vector<std::string>& ids, bool includeTrash, query_options *) "
  /**
   * Retrieves the corresponding documents from the database.
   *
   * @param index - the index from where to get the documents
   * @param collection - the collection from where to get the documents
   * @param ids - the documents id
   * @param includeTrash - include or not the trash
   * @param options - Request options
   * @return json representing the document
   */
  public";

%javamethodmodifiers kuzzleio::Document::mGet(const std::string& index, const std::string& collection, const std::vector<std::string>& ids, bool includeTrash) "
  /**
   * {@link #mGet(String, String, StringVector, Bool, QueryOptions)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::mReplace(const std::string& index, const std::string& collection, const std::string& body, query_options *options) "
  /**
   * Replaces existing documents in the persistent data storage
   *
   * @param index - the index where to replace the document
   * @param collection - the collection where to replace the document
   * @param ids - the documents ids
   * @param body - the content of the document
   * @param options - Request options
   * @return json representing the documents
   */
  public";

%javamethodmodifiers kuzzleio::Document::mReplace(const std::string& index, const std::string& collection, const std::string& body) "
  /**
   * {@link #mReplace(String index, String collection, String id, String body, QueryOptions options)}
   */
  public";

%javamethodmodifiers kuzzleio::Document::mUpdate(const std::string& index, const std::string& collection, const std::string& body, query_options* options) "
  /**
   * Update parts of a documents
   *
   * @param index - the index where to update the document
   * @param collection - the collection where to update the document
   * @param body - json representing documents to update
   * @param options - Request options
   * @return the documents
   */
  public";

%javamethodmodifiers kuzzleio::Document::mUpdate(const std::string& index, const std::string& collection, const std::string& body) "
  /**
   * {@link #mUpdate(String index, String collection, String body, String body)}
   */
  public";
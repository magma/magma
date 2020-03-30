package com.netflix.conductor.contribs.tenancy.server.contribs.tenancy.model;

import java.util.List;

/**
 * @author Viren
 *
 */
public class SearchResult<T> {

	private long totalHits;

	private List<T> results;

	public SearchResult(){

	}

	public SearchResult(long totalHits, List<T> results) {
		super();
		this.totalHits = totalHits;
		this.results = results;
	}

	/**
	 * @return the totalHits
	 */
	public long getTotalHits() {
		return totalHits;
	}

	/**
	 * @return the results
	 */
	public List<T> getResults() {
		return results;
	}

	/**
	 * @param totalHits the totalHits to set
	 */
	public void setTotalHits(long totalHits) {
		this.totalHits = totalHits;
	}

	/**
	 * @param results the results to set
	 */
	public void setResults(List<T> results) {
		this.results = results;
	}


}

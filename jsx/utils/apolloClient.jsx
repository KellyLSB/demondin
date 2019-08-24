// React
import React from 'react';

// Apollo Client Link Setup
import { HttpLink } from 'apollo-link-http';
import { WebSocketLink } from 'apollo-link-ws';
import { ApolloLink, split } from 'apollo-link';
import { onError } from 'apollo-link-error';
import { getMainDefinition } from 'apollo-utilities';

// Apollo Client Libraries
import { ApolloProvider as AP } from "react-apollo";
import { ApolloClient } from 'apollo-client';
import { InMemoryCache } from 'apollo-cache-inmemory';

// GraphQL Backend
export const GRAPHQL_BACKEND = "localhost:8080/graphql";

// Create an HTTP link:
export const httpLink = new HttpLink({
	uri: "http://" + GRAPHQL_BACKEND,
});

// Create a WebSocket link:
export const wsLink = new WebSocketLink({
	uri: "ws://" + GRAPHQL_BACKEND,
	options: { reconnect: true },
});

// Handle link splitting:
export const link = ApolloLink.from([
	onError(({ graphQLErrors, networkError }) => {
		if (graphQLErrors)
			graphQLErrors.map(({ message, locations, path }) =>
				console.log(
					'[GraphQL error]: Message: ', message,
					' Location: ', locations, 
					' Path: ', path,
				),
			);
		if (networkError) 
			console.log('[Network error]: ', networkError);
  }),
	split(
		// split based on operation type
		({ query }) => {
			const { kind, operation } = getMainDefinition(query);
			console.log("Operation:", kind, operation);
			return kind === 'OperationDefinition' && operation === 'subscription';
		},
		wsLink,
		httpLink,
	)
]);

// using the ability to split links, you can send data to each link
// depending on what kind of operation is being sent
export const client = new ApolloClient({ 
	cache: new InMemoryCache(), 
	link: link, 
});

export class ApolloProvider extends React.Component {
	render() {
		return <AP client={client}>
			{this.props.children}
		</AP>;
	}
}

export default ApolloProvider;

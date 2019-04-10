import React from 'react';
import ReactDOM from 'react-dom';

// Apollo
import { ApolloProvider } from "react-apollo";
import { ApolloClient } from 'apollo-client';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { HttpLink } from 'apollo-link-http';
import { onError } from 'apollo-link-error';
import { ApolloLink, split } from 'apollo-link';
import { WebSocketLink } from 'apollo-link-ws';
import { getMainDefinition } from 'apollo-utilities';

const GRAPHQL_BACKEND = "localhost:8080/graphql";

// Create an HTTP link:
const httpLink = new HttpLink({
  uri: "http://" + GRAPHQL_BACKEND,
});

// Create a WebSocket link:
const wsLink = new WebSocketLink({
  uri: "ws://" + GRAPHQL_BACKEND,
  options: { reconnect: true },
});

// Handle link splitting:
const link = ApolloLink.from([
	onError(({ graphQLErrors, networkError }) => {
		if (graphQLErrors)
			graphQLErrors.map(({ message, locations, path }) =>
      	console.log(
					`[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`,
				),
			);
   	if (networkError) console.log(`[Network error]: ${networkError}`);
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
const client = new ApolloClient({ link, cache: new InMemoryCache() });

// Style and Elements
import { Container, Grid, Segment } from 'semantic-ui-react';
import '../../semantic/dist/semantic.min.css';

import Items from './items';
import Cart from './cart';

ReactDOM.render(
  <ApolloProvider client={client}>
    <Container text>
      <Grid columns={2} divided>
        <Grid.Row stretched>
          <Grid.Column width={10}>
            <Items addToCart={(id) => console.log("ID: ", id)} />
          </Grid.Column>
          <Grid.Column width={6}>
						<Segment>
							<Cart />
						</Segment>
            <Segment>Cart Data</Segment>
            <Segment>Checkout</Segment>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Container>
  </ApolloProvider>,
  document.getElementById('root')
);

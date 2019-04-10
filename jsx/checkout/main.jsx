import React from 'react';
import ReactDOM from 'react-dom';

// Apollo
import ApolloClient from "apollo-boost";
import { ApolloProvider } from "react-apollo";

// Client Setup
import { split } from 'apollo-link';
import { HttpLink } from 'apollo-link-http';
import { WebSocketLink } from 'apollo-link-ws';
import { getMainDefinition } from 'apollo-utilities';

// Create an HTTP link:
const httpLink = new HttpLink({
  uri: "http:///graphql",
});

// Create a WebSocket link:
const wsLink = new WebSocketLink({
  uri: "ws:///graphql",
  options: { reconnect: true },
});

// using the ability to split links, you can send data to each link
// depending on what kind of operation is being sent
const client = new ApolloClient({
	link: split(
		// split based on operation type
		({ query }) => {
			const { kind, operation } = getMainDefinition(query);
			return kind === 'OperationDefinition' && operation === 'subscription';
		},
		wsLink,
		httpLink,
	)
});

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
            <Items addToCart={(id) => console.log("ID: ", id)}/>
          </Grid.Column>
          <Grid.Column width={6}>
            <Segment>Cart Data</Segment>
            <Segment>Checkout</Segment>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Container>
  </ApolloProvider>,
  document.getElementById('root')
);

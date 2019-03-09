import React from 'react';
import ReactDOM from 'react-dom';

import ApolloClient from "apollo-boost";
import { ApolloProvider } from "react-apollo";

import { Container, Grid, Segment } from 'semantic-ui-react';
import '../../semantic/dist/semantic.min.css';

import Items from './items';

const client = new ApolloClient();

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

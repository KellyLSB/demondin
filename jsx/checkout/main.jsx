import React from 'react'
import ReactDOM from 'react-dom'
import { Container, Grid, Segment } from 'semantic-ui-react';

import Items from './items'

import '../../semantic/dist/semantic.min.css';

ReactDOM.render(
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
    </Container>,
    document.getElementById('root')
);

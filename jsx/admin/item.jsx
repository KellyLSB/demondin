import React from 'react'
import { Mutation } from "react-apollo";
import gql from "graphql-tag";

import Prices from './prices'

import { Grid, Form, Button } from 'semantic-ui-react'
import FormHelper from '../utils/formHelper'
import GridList from '../utils/gridList'

export default class Item extends FormHelper {
  constructor(props) {
    super(props);
    this.state = {
      item: {}
    };

    if ('item' in props) {
      this.state.item = props.item;
    }

    this.onSubmit = this.onSubmit.bind(this);
  }

  onSubmit(e, updateItem) {
    e.preventDefault();

    console.log(this.state.item);

    updateItem({ variables: {
      input: JSON.stringify(this.state.item),
    } });
  }

  render() {
    return (
      <Grid.Row>
        <Grid.Column>
          <Mutation mutation={gql`
            mutation updateItem($id: ID!, $input: NewItem!) {
              updateItem(id: $id, input: $input) {
                id
                name
                description
                prices {
                  id
                  price
                  afterDate
                  beforeDate
                }
                options {
                  id
                  itemOptionType {
                    id
                    key
                  }
                  values
                }
              }
            }
          `}>
            {(updateItem) => (
              <Form onSubmit={(e) => this.onSubmit(e, updateItem)}>
                <Form.Input label="Name" placeholer="Name"
                  name='name' value={this.state.item.name}
                  onChange={this.onChange} />
                <Form.Field label="Description" placeholder="Description"
                  name='description' control='textarea' rows="3"
                  value={this.state.item.description}
                  onChange={this.onChange} />
                <Button type='submit'>Save</Button>
              </Form>
            )}
          </Mutation>
        </Grid.Column>
        <Grid.Column>
          <Prices item={this.state.item}
            prices={this.state.item.prices} />
        </Grid.Column>
      </Grid.Row>
    );
  }
}

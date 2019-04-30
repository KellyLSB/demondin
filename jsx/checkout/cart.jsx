import React from 'react'
import { Subscription } from "react-apollo";
import gql from "graphql-tag";

import Item from './item'

import { List, Button, Header, Icon, Grid, Segment } from 'semantic-ui-react'

export default class Cart extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Grid.Row>
        <Header attached='top'>
					<Icon name='shopping cart' />
					<Header.Content>Cart</Header.Content>
				</Header>
				<Subscription subscription={gql`
					subscription InvoiceUpdated {
						invoiceUpdated {
							id
							items {
								id
								item {
									name
								}
								itemPrice {
									price
								}
								options {
									itemOptionType {
										id
										key
									}
									values
								}
							}
						}
					}`}>
						{({ data, loading }) => {
							var invoice = data ? data.invoiceUpdated : false;

							//console.log(loading);
							//console.log(invoice);
						
							if (!loading && invoice) return (
								<React.Fragment>
									<Segment attached>
										<Header as='h3'>
											Invoice ID: {invoice.id}
										</Header>
									</Segment>
									{invoice.items.map((item) =>
										<Segment key={item.id} attached>
											<Header as='h4' dividing>{item.item.name}</Header>
											<span>{item.itemPrice.price.toDollars()}</span>
											<List>
												{item.options.map((option) =>
													<List.Item key={option.itemOptionType.id}>
														<List.Header>{option.itemOptionType.key}</List.Header>
														{option.values}
													</List.Item>
												)}
											</List>
										</Segment>											
									)}									
								</React.Fragment>
							);

							return (<Segment attached>Empty Cart</Segment>);
						}	
					}
				</Subscription>
				<Segment attached>
					<Button primary>Checkout</Button>
				</Segment>
			</Grid.Row>
		)
	}
}

import React from 'react'
import { Subscription, Mutation } from "react-apollo";
import gql from "graphql-tag";

import Item from './item'

import { Form, List, Button, Header, Icon, Grid, Segment } from 'semantic-ui-react'

export default class Cart extends React.Component {
  constructor(props) {
    super(props);

		this.onSubmit = this.onSubmit.bind(this);
  }

	onSubmit(e, updateInvoice) {
		e.preventDefault()

		updateInvoice({ variables: {
			input: { 
				submit: true
			}
		} })
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
							subTotal
							demonDin
							taxes
							total
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
											<Header as='h4' dividing>{item.item ? item.item.name : "<Item>"}</Header>
											<span>{item.itemPrice ? item.itemPrice.price.toDollars() : "-.-"}</span>
											<List>
												{item.options.map((option, i) =>
													<List.Item key={option.itemOptionType ? option.itemOptionType.id : `option-${i}`}>
														<List.Header>
															{option.itemOptionType ? option.itemOptionType.key : null}
														</List.Header>
														{option.values}
													</List.Item>
												)}
											</List>
										</Segment>											
									)}
									{invoice.total > 0 ? (
										<Segment attached>
											<List>
												<List.Item>
													<List.Header>SubTotal</List.Header>
													{invoice.subTotal.toDollars()}
												</List.Item>
												<List.Item>
													<List.Header>DemonDin</List.Header>
													{invoice.demonDin.toDollars()}
												</List.Item>
												<List.Item>
													<List.Header>Taxes</List.Header>
													{invoice.taxes.toDollars()}
												</List.Item>
												<List.Item>
													<List.Header>Total</List.Header>
													{invoice.total.toDollars()}
												</List.Item>
											</List>
										</Segment>
									) : null }					
								</React.Fragment>
							);

							return (<Segment attached>Empty Cart</Segment>);
						}	
					}
				</Subscription>
				<Segment attached>
					<Mutation mutation={gql`
						mutation activeInvoice($input: NewInvoice!) {
							activeInvoice(input: $input) {
								id
							}
						}
					`}>
						{(updateInvoice) => (
							<Form onSubmit={(e) => this.onSubmit(e, updateInvoice)}>
								<Button type='submit'>Checkout</Button>
							</Form>
						) }				
					</Mutation>
				</Segment>
			</Grid.Row>
		)
	}
}

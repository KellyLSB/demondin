import React from 'react'
import { Subscription, Mutation } from "react-apollo";
import gql from "graphql-tag";

import Item from './item'

import { Input, Form, List, Button, Header, Icon, Grid, Segment } from 'semantic-ui-react'

import {CardElement, injectStripe} from 'react-stripe-elements';

class Cart extends React.Component {
	constructor(props) {
		super(props);
		
		this.state = {
			values: {},
		};
		
		this.onChange = this.onChange.bind(this);
		this.onSubmit = this.onSubmit.bind(this);
	}
	
	onChange(e, { name, value }) {
		this.setState((state) => {
			state.values[name] = value;
			return state
		} );
	}

	onSubmit(e, updateInvoice) {
		e.preventDefault()
		
		console.log(this.state.values)

		let { token } = this.props.stripe.createToken({
			name: 		this.state.values.cardHolder,
			address_line1: 	this.state.values.cardAddress,
			address_city: 	this.state.values.cardCity,
			address_state: 	this.state.values.cardState,
			address_zip: 	this.state.values.cardZip,
		});

		updateInvoice({ variables: {
			input: { 
				cardToken: token.id,
				submit: true,
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
										<React.Fragment>
											<Segment attached>
												<Grid columns={2}>
													<Grid.Row>
														<Grid.Column>
															<Header sub>SubTotal</Header>
															{invoice.subTotal.toDollars()}
														</Grid.Column>
														<Grid.Column textAlign='right'>
															<Header sub>DemonDin</Header>
															{invoice.demonDin.toDollars()}
														</Grid.Column>
													</Grid.Row>
													<Grid.Row>
														<Grid.Column>
															<Header sub>Taxes</Header>
															{invoice.taxes.toDollars()}
														</Grid.Column>
														<Grid.Column textAlign='right'>
															<Header sub>Total</Header>
															{invoice.total.toDollars()}
														</Grid.Column>
													</Grid.Row>
												</Grid>
											</Segment>
											<Segment attached>
												<Input fluid transparent name="cardHolder" 
															 placeholder="Cardholder"
															 onChange={this.onChange} />
											</Segment>
											<Segment attached>
												<Input fluid transparent name="cardAddress"
															 placeholder="Street Address" 
															 onChange={this.onChange} />
												<Input fluid transparent name="cardCity"
															 placeholder="City" 
															 onChange={this.onChange} />
												<Input fluid transparent name="cardState" 
															 placeholder="State" 
															 onChange={this.onChange} />
												<Input fluid transparent name="cardZip"
															 placeholder="ZIP Code" 
															 onChange={this.onChange} />
											</Segment>
											<Segment attached>
												<CardElement />
											</Segment>
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
										</React.Fragment>
									) : null }		
								</React.Fragment>
							);

							return (<Segment attached>Empty Cart</Segment>);
						}	
					}
				</Subscription>
			</Grid.Row>
		)
	}
}

export default injectStripe(Cart);

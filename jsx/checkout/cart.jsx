import React from 'react'
import { Subscription, Mutation } from "react-apollo";
import gql from "graphql-tag";
import * as EmailValidator from 'email-validator';

import Item from './item'

import { 
	Input, Form, List, Button, Header, 
	Icon, Grid, Segment, Message,
} from 'semantic-ui-react'

import {CardElement, injectStripe} from 'react-stripe-elements';

class Cart extends React.Component {
	constructor(props) {
		super(props);
		
		this.state = {
			errors: {},
			values: {},
		};
		
		this.onError = this.onError.bind(this);
		this.hasError = this.hasError.bind(this);
		this.onChange = this.onChange.bind(this);
		this.onSubmit = this.onSubmit.bind(this);
	}
	
	onError(name, msg = null) {
		this.setState((state) => {
			if (msg === false) {
				delete state.errors[name];
			} else {
				state.errors[name] = msg;
			}
			
			return state;
		} );
	}
	
	hasError(name = null) {
		if (name === null) {
			return Object.keys(this.state.errors).length > 0;
		}
		
		return this.state.errors.hasOwnProperty(name);
	}
	
	onChange(e, { name, value }) {
		// Validate Email Address
		if (name.startsWith("emailAddress")) {
			if ( ! EmailValidator.validate(value)) {
				this.onError(name, `Invalid email address: ${value}`);
				return;
			}
			
			if (name == "emailAddress2") {
				if (value !== this.state.values.emailAddress) {
					this.onError(name, "Please ensure email address authenticity.");
					return;
				}
			}
			
			this.onError(name, false);
		}
	
		this.setState((state) => {
			state.values[name] = value;
			return state;
		} );
	}

	onSubmit(e, updateInvoice) {
		e.preventDefault()
		
		// Clear Stripe Errors
		this.onError("stripe", false);
		
		// Other Errors Block
		if (this.hasError()) return;

		// Create Stripe Token
		this.props.stripe.createToken({
			name: 					this.state.values.cardHolder,
			address_line1: 	this.state.values.cardAddress,
			address_city: 	this.state.values.cardCity,
			address_state: 	this.state.values.cardState,
		}).then(({ error, token }) => {
			if (error) {
				this.onError("stripe", error.message);
				return;
			}
		
			updateInvoice({ variables: {
				input: {
					items: [], // Req?!
					account: { 
						name:  this.state.values.cardHolder,
						email: this.state.values.emailAddress, 
					}, 
					stripeTokenID: token.id,
					submit: true,
				}
			} });
		});
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
							stripeTokenID
							stripeToken
							stripeChargeID
							stripeCharge
							account {
								email
							}
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
											Invoice ID: { invoice.id }
										</Header>
									</Segment>
									{ invoice.items.map((item) =>
										<Segment key={item.id} attached>
											<Header as='h4' dividing>
												{ item.item ? item.item.name : "<Item>" }
											</Header>
											<span>
												{ item.itemPrice ? 
													item.itemPrice.price.toDollars() : "-.-"
												}
											</span>
											<List>
												{ item.options.map((option, i) =>
													<List.Item key={option.itemOptionType ? 
														option.itemOptionType.id : `option-${i}`
													} >
														<List.Header>
															{ option.itemOptionType ? 
																option.itemOptionType.key : null
															}
														</List.Header>
														{ option.values }
													</List.Item>
												) }
											</List>
										</Segment>
									) }
									{ invoice.total > 0 ? (
										<React.Fragment>
											<Segment attached>
												<Grid columns={2}>
													<Grid.Row>
														<Grid.Column>
															<Header sub>SubTotal</Header>
															{ invoice.subTotal.toDollars() }
														</Grid.Column>
														<Grid.Column textAlign='right'>
															<Header sub>DemonDin</Header>
															{ invoice.demonDin.toDollars() }
														</Grid.Column>
													</Grid.Row>
													<Grid.Row>
														<Grid.Column>
															<Header sub>Taxes</Header>
															{ invoice.taxes.toDollars() }
														</Grid.Column>
														<Grid.Column textAlign='right'>
															<Header sub>Total</Header>
															{ invoice.total.toDollars() }
														</Grid.Column>
													</Grid.Row>
												</Grid>
											</Segment>
											{ invoice.stripeToken ? (
												<Segment attached>
													{ invoice.stripeToken.bank_account ? (
														<React.Fragment>
															<Header as='h3'>Bank Account Transfer</Header>
															{ console.log(invoice.stripeToken.bank_account) }
														</React.Fragment>
													) : null }
													{ invoice.stripeToken.card ? (
														<React.Fragment>
															<Header as='h3'>Card Applied</Header>
															<span>
																{     invoice.stripeToken.card.brand }
																 x{   invoice.stripeToken.card.last4 }
																 :: { invoice.stripeToken.card.name  }
															</span>
															<br />
															<span>
																{ invoice.stripeToken.used ? 
																	'Charged' : 'Uncharged' 
																}
															</span>
														</React.Fragment>
													) : null }
												</Segment>
											) : (
												<React.Fragment>
													<Segment attached>
														<Input fluid transparent 
																	 error={ this.hasError("emailAddress") }
																	 name="emailAddress" 
																	 placeholder="Email Address" 
																	 onChange={this.onChange} />
														<Input fluid transparent 
																	 error={ this.hasError("emailAddress2") }
																	 name="emailAddress2" 
																	 placeholder="Confirm Email Address" 
																	 onChange={this.onChange} />
													</Segment>
													<Segment attached>
														<Input fluid transparent
																	 error={ this.hasError("cardHolder") } 
																	 name="cardHolder" 
																	 placeholder="Cardholder"
																	 onChange={this.onChange} />
													</Segment>
													<Segment attached>
														<Input fluid transparent 
																	 error={ this.hasError("cardAddress") }
																	 name="cardAddress"
																	 placeholder="Street Address" 
																	 onChange={this.onChange} />
														<Input fluid transparent 
																	 error={ this.hasError("cardCity") }
																	 name="cardCity"
																	 placeholder="City" 
																	 onChange={this.onChange} />
														<Input fluid transparent 
																	 error={ this.hasError("cardState") }
																	 name="cardState" 
																	 placeholder="State" 
																	 onChange={this.onChange} />
													</Segment>
													<Segment attached>
														<CardElement />
													</Segment>
												</React.Fragment>
											) }
											{ invoice.stripeCharge ? (
												<Segment attached>
													<Header as='h3'>Charged and Checked Out</Header>
													<span>
														{ invoice.account ? invoice.account.email : "<email>" } 
														:: { invoice.stripeCharge.paid ? 
															"Paid" : `Status: ${invoice.stripeCharge.status}`
														}
													</span>
												</Segment>
											) : (
												<Segment attached>
													<Mutation mutation={gql`
														mutation activeInvoice($input: NewInvoice!) {
															activeInvoice(input: $input) {
																id
															}
														}
													`}>
														{ (updateInvoice) => (
															<Form error={this.hasError()} 
																onSubmit={ (e) => this.onSubmit(e, updateInvoice) }>
																<Message error header='Error checking out'
																	list={Object.values(this.state.errors)}
																/>
																<Button type='submit'>Checkout</Button>
															</Form>
														) }				
													</Mutation>
												</Segment>
											) }
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

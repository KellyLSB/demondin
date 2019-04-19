import React from 'react'
import { Subscription } from "react-apollo";
import gql from "graphql-tag";

import Item from './item'

import { Header, Grid } from 'semantic-ui-react'

export default class Cart extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Grid.Row>
        <Header>Badges</Header>
					<Subscription subscription={gql`
					subscription InvoiceUpdated {
						invoiceUpdated {
							id
							items {
								item {
									name
								}
								itemPrice {
									price
								}
								options {
        					optionType {
										key
        					}
        					values
      					}
							}
						}
	 				}`}>
						{({ data, loading }) => {
							var invoice = data ? data.invoiceUpdated : false;

							console.log(loading);
							console.log(invoice);
							
							if (!loading && invoice) return (
								<React.Fragment>
									{invoice.id}
									{invoice.items.map((item) =>
										<div key={item.id} style={{border: "underline #000 solid"}}>
											<h3>{item.name}</h3>
											<span>{item.itemPrice.price}</span>
											{item.options.map((option) =>
												<span>{option.optionType.key}: {option.values}</span>
											)}
										</div>
									)}
								</React.Fragment>
    					);

							return (<span>Loading</span>);
						}	}
					</Subscription>
 	     </Grid.Row>
 	   )
  }
}

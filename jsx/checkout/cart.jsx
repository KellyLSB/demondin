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
								id
								itemID
								itemPriceID
							}
						}
	 				}`}>
						{({ data, loading }) => {
							var invoice = data ? data.invoiceUpdated : false;
							console.log(invoice);
							
							if (!loading && invoice) return (
								<React.Fragment>
									{invoice.id}
									{invoice.items.map((item) =>
										<span key={item.id} style={{border: "underline #000 solid"}}>{item.id}</span>
									)}
								</React.Fragment>
    					)

							return (<span>Loading</span>);
						}	}
					</Subscription>
 	     </Grid.Row>
 	   )
  }
}

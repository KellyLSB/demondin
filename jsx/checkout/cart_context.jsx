import React from 'react'

// Likely will use this to pass event callbacks down through the stack.
export const CartContext = React.createContext({
  // Handle Invoice Creation
  cartContents: [],
  // Request Completion of Options
  addToCart: (id) => () =>
    fetch("/shop/invoicing/addToCart.json", {
      method: "POST",
      body: JSON.stringify(id)
    }).then((response) =>
      response.json()
    ).then((data) => {
      console.log("Need to link this process with the cart context (accept options)")
      console.log(data);
      console.log(this);
    })
})
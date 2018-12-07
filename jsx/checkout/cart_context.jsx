import React from 'react'

// Likely will use this to pass event callbacks down through the stack.
export default const CartContext = React.createContext({
  // Handle Invoice Creation
  cartContents: [],
  // Request Completion of Options
  addToCart: (id, options) =>
    return null
})
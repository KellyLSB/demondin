import React from 'react'
import { CartContext } from './cart_context'

// Likely will use this to pass event callbacks down through the stack.
export default class App extends React.Component {
  
  
  render() {
    return (
      <CartContext.Provider>
      </CartContext.Provider>
    )
  }
}
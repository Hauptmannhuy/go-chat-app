import { createContext } from "react";

export const SearchContext = createContext({
  handleSearchQuery: null, 
  setEmptyInput: null
})
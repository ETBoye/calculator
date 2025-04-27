

export type HistoryResponse = {
    self: string,
    first: string,
    prev?: string,
    next?: string,
    last: string,
    items: HistoryRow[]
    error: string,
}
  
  
export type ComputeResponse = {
    error?: string,
    historyRow: HistoryRow
}
  

export type CalculationResult = {
    input: string,
    result?: RationalNumber
    errorId?: string,
}
  
export type RationalNumber =  {
    num: string,
    denom: string, 
    estimate: string
}
  
export type HistoryRow = {
    calculationId: number,
    calculation: CalculationResult,
}
  
export type HistoryData = {
    nextUrl?: string,
    items: HistoryRow[]
}

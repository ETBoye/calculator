'use client';

import { KeyboardEventHandler, useEffect, useState } from "react";
import { v4 } from "uuid"

const SESSION_ID_KEY = "session-id";

type CalculationResult = {
  input: string,
  result?: RationalNumber
  errorId?: string,
}

type RationalNumber =  {
  num: string,
  denom: string, 
  estimate: string
}

type HistoryRow = {
  calculationId: number,
  calculation: CalculationResult,
}

type HistoryData = {
  nextUrl?: string,
  items: HistoryRow[]
}

type ComputeResponse = {
  error?: string,
  historyRow: HistoryRow
}


const API_BASE = "/api";

type HistoryResponse = {
  self: string,
  first: string,
  prev?: string,
  next?: string,
  last: string,
  items: HistoryRow[]
  error: string,
}

function Result({calculationResult} : {calculationResult: CalculationResult}) {
  return (
    <>
      {
        calculationResult.errorId ? <p>{calculationResult.errorId}</p>
        : <p>{calculationResult.result!.num}/{calculationResult.result!.denom} = {calculationResult.result!.estimate}</p>
      }
    </>
  )
}

function SessionIdTool({onChangeSessionId, sessionId}: 
  {onChangeSessionId: (newSessionId: string) => void, sessionId: string}) {
  const [sessionIdInputValue, setSessionIdInputValue] = useState<string>("");


  return (
    <>
    <p> Current session id: {sessionId} </p>
    <div>
      <p> Change session id: </p>
      <input value={sessionIdInputValue} onChange={ e => setSessionIdInputValue(e.currentTarget.value)}></input>
      <button onClick={() => onChangeSessionId(sessionIdInputValue)}>Change</button>
    </div>
    </>
  )
}

function HistorySection({historyData, onFetchNextHistoryPage}: 
  {historyData: HistoryData, onFetchNextHistoryPage: () => Promise<void>}) {
  return (
    <>
    <table>
          
          <thead>
          <tr>
            <th>Input</th>
            <th>Result</th>
          </tr>
          </thead>
          <tbody>
          {
            historyData.items.map( item => 
              <tr key={item.calculationId}>
                <td>{item.calculation.input}</td>

                <td>
                  <Result calculationResult={item.calculation}></Result>
                </td>
              </tr>
            )
          }
          </tbody>
        </table>
        {
          historyData.nextUrl ? <button onClick={() => onFetchNextHistoryPage()}>Fetch more..</button> : null
        }
    </>
  )
}

function InputSection({onCalculate, currentResult}: 
  {onCalculate: (input: string) => Promise<void>, currentResult?: CalculationResult | null}) {
  const [calculationInputValue, setCalculationInputValue] = useState<string>("");

  function onKeyUp(e: any) {
    if (e && e.key == 'Enter') {
      onCalculate(calculationInputValue)
    }
  }

  return (
    <>
     <div style={{marginTop: '100px'}}>
        
        <input placeholder="1+(3+4)" 
          value={calculationInputValue} 
          onChange={ e => setCalculationInputValue(e.currentTarget.value)} 
          onKeyUp={(e) => onKeyUp(e)}></input>
        <button onClick={() => onCalculate(calculationInputValue)}>Calculate!</button>
        { currentResult ? <Result calculationResult={currentResult}></Result> : null}
        
    </div>
    </>
  )
}


async function fetchWithErrorHandle<T>(
  input: string | URL | globalThis.Request,
  init?: RequestInit,
): Promise<T> {
  try {
    const response = await fetch(input, init)

    if (!response.ok) {
      alert(`Received response status: ${response.status} with body ${JSON.stringify(await response.json())}`)
      throw new Error(`Response not ok, status: ${response.status}`)
    }

    return await response.json()

  } catch(e) {
    alert("Received error from backend!");
    throw e;
  }
}

function sessionIdIsValid(sessionId: string): boolean {
  return /^[a-zA-Z0-9\-]*$/.test(sessionId) && sessionId.length > 1 && sessionId.length < 100
}

export default function Home() {
  const [sessionId, setSessionId] = useState<string|null>(null);
  const [historyData, setHistoryData] = useState<HistoryData>({items: []})
  const [currentResult, setCurrentResult] = useState<CalculationResult | null>();

  useEffect(() => {
    let sessionId = window.localStorage.getItem(SESSION_ID_KEY) ?? v4()

    if (!sessionIdIsValid(sessionId)) {
      sessionId = v4()
    }
    changeSessionId(sessionId)
  }, [])
  
  useEffect(() => {
    if (sessionId != null) {
      fetchHistoryForNewSessionId();
    }
  }, [sessionId])

  async function fetchHistoryForNewSessionId() {
    const historyResponse: HistoryResponse = await fetchWithErrorHandle(`/api/sessions/${sessionId}/history`)
    
    setHistoryData({
      nextUrl: historyResponse.next,
      items: historyResponse.items
    })
  }

  async function fetchNextHistoryPage() {
    if (historyData.nextUrl) {
      const historyResponse: HistoryResponse = await fetchWithErrorHandle(`${API_BASE}${historyData.nextUrl}`)
      
      setHistoryData((previousHistoryData) => {
        return {
          nextUrl: historyResponse.next,
          items: previousHistoryData.items.concat(historyResponse.items)
        }
      });
    }
  }

  async function changeSessionId(sessionId: string) {
    window.localStorage.setItem(SESSION_ID_KEY, sessionId)
    setSessionId(sessionId)
  }

  async function calculate(input: string) {

    const ComputeResponse: ComputeResponse = await fetchWithErrorHandle(`/api/sessions/${sessionId}/compute`, {
        method: 'POST',
        headers: { "Content-Type": "application/json"},
        body: JSON.stringify({
          "input": input
        })
      })

    setCurrentResult(ComputeResponse.historyRow.calculation)
    await fetchHistoryForNewSessionId()
  }

  

  return (
    <>
        <h1>Calculator</h1>

        
        { sessionId ?         <SessionIdTool onChangeSessionId={changeSessionId} sessionId={sessionId}></SessionIdTool> : null}
        
        <InputSection currentResult={currentResult} onCalculate={calculate}></InputSection>

    
        <HistorySection historyData={historyData} onFetchNextHistoryPage={fetchNextHistoryPage}></HistorySection>
        
        

    </>
  );
}

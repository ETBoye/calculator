import { HistoryData } from "../model";

import styles from "../page.module.css"
import { Result } from "./result";


export function HistoryCard({historyData, onFetchNextHistoryPage, sessionId}: 
    {historyData: HistoryData, onFetchNextHistoryPage: () => Promise<void>, sessionId: string}) {
    return (
      <div className={`${styles.historySection} ${styles.gridCard}`}>
        <h2>History</h2>
        <p>Session id: {sessionId}</p>
        <table className={styles.historyTable}>
              
          <thead>
          <tr>
            <th>Input</th>
            <th>Output</th>
          </tr>
          </thead>
          <tbody>
          {
            historyData.items.map( item => 
              <tr key={item.calculationId} className={styles.historyRow}>
                <td style={{textAlign: 'center'}}>{item.calculation.input}</td>
  
                <td style={{textAlign: 'center'}}>
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
      </div>
    )
}
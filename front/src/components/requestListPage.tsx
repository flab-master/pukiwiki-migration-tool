import './requestPage.css'

type MigrationRequest = {
  id: number
  userID: string
  pageID: string
  requestedAt: string
  status: '申請中' | '完了' | '失敗'
}

const mockRequests: MigrationRequest[] = [
  {
    id: 1,
    userID: 'fujimoto2025',
    pageID: 'notion-page-8f42a1',
    requestedAt: '2026-05-01 18:20',
    status: '申請中',
  },
  {
    id: 2,
    userID: 'tanaka2024',
    pageID: 'notion-page-41c9d0',
    requestedAt: '2026-04-30 14:05',
    status: '完了',
  },
  {
    id: 3,
    userID: 'suzuki-seminar',
    pageID: 'notion-page-73aa19',
    requestedAt: '2026-04-29 09:48',
    status: '失敗',
  },
]

const statusClassName = (status: MigrationRequest['status']) => {
  switch (status) {
    case '完了':
      return 'status-badge done'
    case '失敗':
      return 'status-badge failed'
    default:
      return 'status-badge pending'
  }
}

export default function RequestListPage() {
  return (
    <main className="list-page">
      <section className="panel list-panel">
        <div className="page-heading list-heading">
          <div>
            <p className="eyebrow">Requests</p>
            <h1>移行申請一覧</h1>
            <p>バックエンド接続前のため、現在はモックデータを表示しています。</p>
          </div>
          <span className="count-badge">{mockRequests.length} 件</span>
        </div>

        <div className="request-table-wrap">
          <table className="request-table">
            <thead>
              <tr>
                <th>状態</th>
                <th>userID</th>
                <th>pageID</th>
                <th>申請日時</th>
              </tr>
            </thead>
            <tbody>
              {mockRequests.map((request) => (
                <tr key={request.id}>
                  <td>
                    <span className={statusClassName(request.status)}>{request.status}</span>
                  </td>
                  <td>{request.userID}</td>
                  <td>{request.pageID}</td>
                  <td>{request.requestedAt}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>
    </main>
  )
}

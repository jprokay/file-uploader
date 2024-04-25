import GetStartedAlert from "@/components/alerts/get-started-alert";
import { ContactsDataTable } from "@/components/tables/contacts-data-table";
import { ContactsWithTotal, client } from "@/lib/contacts-api";
import { cookies } from 'next/headers'

type GetDataProps = {
  limit: number
  offset: number
}
async function getData(props: GetDataProps): Promise<ContactsWithTotal> {

  const res = await client.GET("/contacts", {
    params: {
      query: props,
      cookie: { userId: cookies().get("userId")?.value || "" }
    },
    headers: { Cookie: cookies().toString() },
  })


  if (res.error) {
    return { contacts: [], total: 0 }
  }

  return res.data
}

export default async function Contacts() {
  const data = await getData({ limit: 100, offset: 1 })
  return (
    <main className="container mx-auto py-10">
      <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">
        Contacts
      </h1>
      {data.total < 1 && <div className="my-4">
        <GetStartedAlert />
      </div>}
      <ContactsDataTable data={data.contacts} rowCount={data.total} />
    </main>
  )
}

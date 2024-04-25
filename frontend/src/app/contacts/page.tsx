import GetStartedAlert from "@/components/alerts/get-started-alert";
import { ContactsDataTable } from "@/components/tables/contacts-data-table";

export default async function Contacts() {
  return (
    <main className="container mx-auto py-10">
      <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">
        Contacts
      </h1>
      {<div className="my-4">
        <GetStartedAlert />
      </div>}
      <ContactsDataTable />
    </main>
  )
}

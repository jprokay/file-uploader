import GetStartedAlert from "@/components/alerts/get-started-alert";
import { DirectoriesDataTable } from "@/components/tables/directories-data-table";

export default async function Home() {

  return (
    <div>
      <main className="container mx-auto py-10">
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">
          Directories
        </h1>
        <div className="my-4">
          <GetStartedAlert />
        </div>
        <DirectoriesDataTable />
      </main>
    </div>
  );
}

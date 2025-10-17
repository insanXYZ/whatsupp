import { useMutation, useQuery } from "@tanstack/react-query";
import { API } from "./axios";
import { error } from "console";
import { AxiosError } from "axios";
import { ResponseSchema } from "@/app/dto";

function Query(queryKey: any[], url: string) {
  useQuery({
    queryKey,
    queryFn: async () => {
      const { data } = await API.get(url);
      return data;
    },
  });
}

export enum HttpMethod {
  POST = "POST",
  PUT = "PUT",
  DELETE = "DELETE",
}

interface Mutate {
  url: string;
  body: any;
  method: HttpMethod;
}

function Mutation<T = any>(mutationKey: any[]) {
  return useMutation<
    ResponseSchema<T>,
    AxiosError<ResponseSchema, any>,
    Mutate
  >({
    mutationFn: async ({ url, body, method }: Mutate) => {
      const res = await API({
        url: url,
        data: body,
        method: method,
      });

      return res.data;
    },
    mutationKey,
  });
}

export { Query, Mutation };

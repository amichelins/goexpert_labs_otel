# goexpert_labs_otel

O projeto se baseia em dois sistemas.

-- servico_req
   Ele recebe o cep e o valida.
   Se o cep for valido, ele chama ao segundo serviço servico_orc
   Se for invalido envia uma resposta negativa.

-- servico_orc
   Recebe o cep do servico_req
   Valida a Weather Api Key seja fornecida e que o cep seja valido
   Se as validações forem bem sucedidas, ele chama a ViaCep para obter o nome da cidade
   No caso de erro envia uma resposta negativa ao servico_req,
   no caso de obter o nome da cidade ele chama a Weather Api para obter os dados da temperatura
   Se obtemos uma resposta positiva, usamos fastjson para obter os dados da tempertura.

   Os dados são preparados e enviados de volta ao servico_req

Estrutura:

  docker-compose.yaml Configuração dos containeres

  Pastas:
    docker   Configuração do otel-collector
    servicos Fontes dos projetos
        servico_orc Serviço que consulta os dados da Cidade e da temperatura
        servico_req Serviço que recebe o cep valida e repassa para o proximo serviço

Chamada:
    No projeto temos as pastas /servicos/servico_orc/api e /servicos/servico_req/api
    elas tem o arquivo http.http com chamadas prontas. Uma com um cep errado e outra com cep correto.


Build:
    docker compose up --build -d -V


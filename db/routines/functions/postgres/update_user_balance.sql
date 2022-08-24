create or replace function update_user_balance_v1_1(user_mobile character, credit_amount integer, transaction_id integer, external_log_id integer) returns void
    language plpgsql
as
$$
declare
__version int;
__affected_rows int := 0;
__retry_count int := 0;
begin

        while __retry_count < 3 loop
select _version
            into __version
from users where mobile = user_mobile
    limit 1;

if __version is null then
                raise EXCEPTION 'user not found. mobile: "%"', user_mobile using DETAIL = '{"code":1}';
end if;

update users
set balance = users.balance + credit_amount,
    _version = _version + 1
where mobile = user_mobile
  and users._version = __version;

GET DIAGNOSTICS __affected_rows := ROW_COUNT;

if __affected_rows != 0 then

begin
update transactions
set status = 2, usage_log_id = external_log_id, amount = credit_amount
where id = transaction_id;

GET DIAGNOSTICS __affected_rows := ROW_COUNT;
if __affected_rows = 0 then
                        raise EXCEPTION 'no transaction found to update';
end if;

                    return;
exception
                    when others then
                        raise EXCEPTION 'unknown error when updating transaction status: %s', SQLERRM using DETAIL = '{"code":2}';
end;

else
                __retry_count := __retry_count + 1;
end if;
end loop;

        raise EXCEPTION 'unknown error to charge wallet by credit code' using DETAIL = '{"code":3}';

end;
$$;